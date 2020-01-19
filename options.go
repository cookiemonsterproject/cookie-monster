package cookiemonster

import (
	"os"
	"runtime"
	"time"
)

const defaultBackoff = 10 * time.Second

type DigesterOptionFunc func(*Digester)

func SetWorkers(workers int) DigesterOptionFunc {
	return func(d *Digester) {
		d.infoF("setting workers to %d", workers)
		d.workers = workers
		d.workChan = make(chan []Cookie, workers)
	}
}

func SetBackoff(backoff Backoff) DigesterOptionFunc {
	return func(d *Digester) {
		d.infoF("setting backoff")
		d.backoff = backoff
	}
}

func SetInfoLog(logger Logger) DigesterOptionFunc {
	return func(d *Digester) {
		d.infoF("setting info logger")
		d.infoLogger = logger
	}
}

func SetErrorLog(logger Logger) DigesterOptionFunc {
	return func(d *Digester) {
		d.infoF("setting error logger")
		d.errorLogger = logger
	}
}

func SetStopSignals(signals ...os.Signal) DigesterOptionFunc {
	return func(d *Digester) {
		d.infoF("setting %d stop signals", len(signals))
		d.stopSignals = signals
	}
}

func (d *Digester) handleDefaults() {
	if d.workers == 0 {
		w := runtime.NumCPU()
		d.infoF("workers not set, defaulting to %d", w)
		d.workers = w
		d.workChan = make(chan []Cookie, w)
	}

	if d.backoff == nil {
		d.infoF("backoff not set, defaulting to ConstantBackoff(10 * time.Second)")
		d.backoff = ConstantBackoff(defaultBackoff)
	}
}
