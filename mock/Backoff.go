package mock

import (
	"time"

	"github.com/cookiejars/cookiejar"
)

var _ cookiejar.Backoff = &Backoff{}

type Backoff struct {
	NextFn      func()
	NextInvoked bool

	CurrentFn      func() time.Duration
	CurrentInvoked bool

	ResetFn      func()
	ResetInvoked bool
}

func (b *Backoff) Next() {
	b.NextInvoked = true
	b.NextFn()
}

func (b *Backoff) Current() time.Duration {
	b.CurrentInvoked = true
	return b.CurrentFn()
}

func (b *Backoff) Reset() {
	b.ResetInvoked = true
	b.ResetFn()
}
