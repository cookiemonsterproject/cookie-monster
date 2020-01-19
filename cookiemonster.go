package cookiemonster

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"plugin"
	"sync"
	"sync/atomic"
	"time"
)

type Cookie interface {
	ID() string
	Content() interface{}
	Metadata() map[string]string
}

type Jar interface {
	Retrieve() ([]Cookie, error)
	Retire(Cookie) error
}

type DigestFn func(Cookie) error

type Digester interface {
	Start(fn DigestFn) error
	Stop()
}

type digester struct {
	workers        int
	workChan       chan []Cookie
	infoLogger     Logger
	errorLogger    Logger
	stopSignals    []os.Signal
	jar            Jar
	backoff        Backoff
	running        atomic.Value
	workersWG      sync.WaitGroup
	orchestratorWG sync.WaitGroup
	mux            sync.Mutex
}

func NewDigesterWithPlugin(jarPath string, options ...DigesterOptionFunc) (Digester, error) {
	plug, err := plugin.Open(jarPath)
	if err != nil {
		return nil, err
	}

	sym, err := plug.Lookup("Jar")
	if err != nil {
		return nil, err
	}

	jar, ok := sym.(Jar)
	if !ok {
		return nil, errors.New("unexpected type from module symbol")
	}

	return NewDigester(jar, options...), nil
}

func NewDigester(jar Jar, options ...DigesterOptionFunc) Digester {
	d := &digester{jar: jar}
	d.running.Store(false)

	for _, option := range options {
		option(d)
	}
	d.infoF("handling defaults")
	d.handleDefaults()

	return d
}

func (d *digester) Start(fn DigestFn) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	if d.isRunning() {
		return errors.New("digester is already running")
	}

	d.infoF("starting digester")
	d.running.Store(true)
	d.startWorkers(fn)
	d.startOrchestrator()

	if len(d.stopSignals) > 0 {
		d.infoF("waiting for OS signals...")
		d.waitForSignals(d.stopSignals...)
		d.Stop()
	}

	return nil
}

func (d *digester) Stop() {
	d.infoF("stopping digester")
	d.running.Store(false)
	d.orchestratorWG.Wait()
	close(d.workChan)
	d.workersWG.Wait()
}

func (d *digester) startWorkers(fn DigestFn) {
	d.workersWG.Add(d.workers)

	work := func(workerID int) {
		defer d.infoF("worker %d: stopping", workerID)
		defer d.workersWG.Done()

		for cc := range d.workChan {
			d.infoF("worker %d: handling %d messages", workerID, len(cc))
			for _, c := range cc {
				d.infoF("worker %d: digesting message %s", workerID, c.ID())
				if err := fn(c); err != nil {
					d.errorF("worker %d: could not digest message %s: %s", workerID, c.ID(), err)
					continue
				}

				d.infoF("worker %d: retiring message %s", workerID, c.ID())
				if err := d.jar.Retire(c); err != nil {
					d.errorF("worker %d: could not retire message %s: %s", workerID, c.ID(), err)
					continue
				}
			}
		}
	}

	d.infoF("starting %d workers", d.workers)
	for i := 0; i < d.workers; i++ {
		workerID := i + 1
		d.infoF("starting worker %d", workerID)
		go work(workerID)
	}
}

func (d *digester) startOrchestrator() {
	d.orchestratorWG.Add(1)

	orchestrate := func() {
		defer d.infoF("orchestrator: stopping")
		defer d.orchestratorWG.Done()

		for {
			if !d.isRunning() {
				break
			}

			currBackoff := d.backoff.Current()
			d.infoF("orchestrator: sleeping for %s", currBackoff.String())
			time.Sleep(currBackoff)

			d.infoF("orchestrator: retrieving cookies from jar")
			cc, err := d.jar.Retrieve()
			if err != nil {
				d.errorF("orchestrator: failed to retrieve from jar: %s", err)
				continue
			}

			if len(cc) == 0 {
				d.infoF("orchestrator: found an empty jar")
				d.backoff.Next()
				continue
			}

			d.backoff.Reset()
			d.infoF("orchestrator: digesting %d cookies", len(cc))
			d.workChan <- cc
		}
	}

	d.infoF("starting orchestrator")
	go orchestrate()
}

func (d *digester) waitForSignals(signals ...os.Signal) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, signals...)
	c := <-signalChan
	d.infoF("signal %s triggered", c)
}

func (d *digester) isRunning() bool {
	return d.running.Load().(bool)
}

func (d *digester) infoF(format string, args ...interface{}) {
	if d.infoLogger != nil {
		d.infoLogger.Printf(fmt.Sprintf("[INFO] %s", format), args...)
	}
}

func (d *digester) errorF(format string, args ...interface{}) {
	if d.errorLogger != nil {
		d.errorLogger.Printf(fmt.Sprintf("[ERROR] %s", format), args...)
	}
}
