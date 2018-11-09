package cookiejar

import (
	"os"
	"runtime"
	"time"
)

const defaultBackoff = 10 * time.Second

type DigesterOptionFunc func(*digester)

func SetWorkers(workers int) DigesterOptionFunc {
	return func(d *digester) {
		d.workers = workers
		d.workChan = make(chan []Cookie, workers)
	}
}

func SetBackoff(backoff Backoff) DigesterOptionFunc {
	return func(d *digester) {
		d.backoff = backoff
	}
}

func SetInfoLog(logger Logger) DigesterOptionFunc {
	return func(d *digester) {
		d.infoLogger = logger
	}
}

func SetErrorLog(logger Logger) DigesterOptionFunc {
	return func(d *digester) {
		d.errorLogger = logger
	}
}

func SetStopSignals(signals ...os.Signal) DigesterOptionFunc {
	return func(d *digester) {
		d.stopSignals = signals
	}
}

func (d *digester) handleDefaults() {
	if d.workers == 0 {
		w := runtime.NumCPU()
		d.workers = w
		d.workChan = make(chan []Cookie, w)
	}

	if d.backoff == nil {
		d.backoff = ConstantBackoff(defaultBackoff)
	}
}
