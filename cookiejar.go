package cookiejar

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"time"
)

const errAlreadyRunning = "digester is already running"

type DigestFn func(cookie Cookie) error

type Digester interface {
	Start(fn DigestFn, signals ...os.Signal) error
}

type digester struct {
	workers  int
	workChan chan Cookie
	jar      Jar
	backoff  Backoff
	running  bool
	wg       sync.WaitGroup
}

func NewDigester(workers int, jar Jar, backoff Backoff) Digester {
	return &digester{
		workers:  workers,
		workChan: make(chan Cookie, workers),
		jar:      jar,
		backoff:  backoff,
	}
}

func (d *digester) Start(fn DigestFn, signals ...os.Signal) error {
	if d.running {
		return errors.New(errAlreadyRunning)
	}

	d.running = true
	d.startWorkers(fn)
	d.startOrchestrator()

	if len(signals) > 0 {
		d.waitForSignal(signals...)
	}

	d.stop()

	return nil
}

func (d *digester) startWorkers(fn DigestFn) {
	d.wg.Add(d.workers)

	work := func() {
		defer d.wg.Done()

		for c := range d.workChan {
			if err := fn(c); err != nil {
				// todo: send to error channel

				continue
			}

			if err := c.Delete(); err != nil {
				// todo: send to error channel
			}
		}
	}
	for i := 0; i < d.workers; i++ {
		go work()
	}
}

func (d *digester) startOrchestrator() {
	d.wg.Add(1)

	orchestrate := func() {
		defer d.wg.Done()

		for d.running {
			select {
			case <-time.After(d.backoff.Current()):
				cc, err := d.jar.Fetch()
				if err != nil {
					// todo: send to error channel

					continue
				}

				if cc == nil {
					d.backoff.Next()

					continue
				}

				d.backoff.Reset()

				for _, c := range cc {
					d.workChan <- c
				}
			}
		}
	}

	go orchestrate()
}

func (d *digester) waitForSignal(signals ...os.Signal) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)
	<-c
}

func (d *digester) stop() {
	d.running = false
	close(d.workChan)
	d.wg.Wait()
}

type Cookie interface {
	Content() ([]byte, error)
	Delete() error
}

type Jar interface {
	Fetch() ([]Cookie, error)
}
