package cookiejar

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

type Cookie interface {
	ID() string
	Content() (interface{}, error)
}

type Jar interface {
	Retrieve() ([]Cookie, error)
	Retire(cookie Cookie) error
}

type DigestFn func(cookie Cookie) error

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

func NewDigester(jar Jar, options ...DigesterOptionFunc) Digester {
	d := &digester{jar: jar}
	d.running.Store(false)

	for _, option := range options {
		option(d)
	}
	d.handleDefaults()

	return d
}

func (d *digester) Start(fn DigestFn) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	if d.isRunning() {
		return errors.New("digester is already running")
	}

	d.running.Store(true)
	d.startWorkers(fn)
	d.startOrchestrator()

	if len(d.stopSignals) > 0 {
		d.waitForSignals(d.stopSignals...)
		d.Stop()
	}

	return nil
}

func (d *digester) Stop() {
	d.mux.Lock()
	defer d.mux.Unlock()

	d.running.Store(false)
	d.orchestratorWG.Wait()
	close(d.workChan)
	d.workersWG.Wait()
}

func (d *digester) startWorkers(fn DigestFn) {
	d.workersWG.Add(d.workers)

	work := func() {
		defer d.workersWG.Done()

		for cc := range d.workChan {
			for _, c := range cc {
				d.infoF("digesting message %s", c.ID())
				if err := fn(c); err != nil {
					d.errorF("could not digest message %s: %s", c.ID(), err)
					continue
				}

				if err := d.jar.Retire(c); err != nil {
					d.errorF("could not retire message %s: %s", c.ID(), err)
					continue
				}
			}
		}
	}

	for i := 0; i < d.workers; i++ {
		go work()
	}
}

func (d *digester) startOrchestrator() {
	d.orchestratorWG.Add(1)

	orchestrate := func() {
		defer d.orchestratorWG.Done()

		for {
			if !d.isRunning() {
				break
			}

			time.Sleep(d.backoff.Current())

			cc, err := d.jar.Retrieve()
			if err != nil {
				// todo: send to error channel

				continue
			}

			if len(cc) == 0 {
				d.backoff.Next()

				continue
			}

			d.backoff.Reset()

			d.workChan <- cc
		}
	}

	go orchestrate()
}

func (d *digester) waitForSignals(signals ...os.Signal) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)
	<-c
}

func (d *digester) isRunning() bool {
	return d.running.Load().(bool)
}

func (d *digester) infoF(format string, args ...interface{}) {
	if d.infoLogger != nil {
		d.infoLogger.Printf(format, args...)
	}
}

func (d *digester) errorF(format string, args ...interface{}) {
	if d.errorLogger != nil {
		d.errorLogger.Printf(format, args...)
	}
}
