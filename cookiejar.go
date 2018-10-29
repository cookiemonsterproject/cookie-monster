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
	Content() (interface{}, error)
	Done() error
}

type Jar interface {
	Retrieve() ([]Cookie, error)
}

type DigestFn func(cookie Cookie) error

type Digester interface {
	Start(fn DigestFn, signals ...os.Signal) error
	Stop()
}

type digester struct {
	workers  int
	workChan chan []Cookie

	jar     Jar
	backoff Backoff

	running        atomic.Value
	workersWG      sync.WaitGroup
	orchestratorWG sync.WaitGroup
	mux            sync.Mutex
}

func NewDigester(workers int, jar Jar, backoff Backoff) Digester {
	d := &digester{
		workers:  workers,
		workChan: make(chan []Cookie, workers),
		jar:      jar,
		backoff:  backoff,
	}
	d.running.Store(false)

	return d
}

func (d *digester) Start(fn DigestFn, signals ...os.Signal) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	if d.isRunning() {
		return errors.New("digester is already running")
	}

	d.running.Store(true)
	d.startWorkers(fn)
	d.startOrchestrator()

	if len(signals) > 0 {
		d.waitForSignals(signals...)
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
				if err := fn(c); err != nil {
					// todo: send to error channel

					continue
				}

				if err := c.Done(); err != nil {
					// todo: send to error channel

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
