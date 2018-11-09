package cookiejar

import (
	"sync"
	"time"
)

type Backoff interface {
	Next()
	Current() time.Duration
	Reset()
}

func ConstantBackoff(amount time.Duration) Backoff {
	return &backoff{intervals: []time.Duration{0, amount}}
}

func ExponentialBackoff(n int, initialAmount time.Duration) Backoff {
	ret := make([]time.Duration, n+1)
	next := initialAmount
	for i := range ret {
		if i == 0 {
			ret[0] = 0
			continue
		}
		ret[i] = next
		next *= 2
	}

	return &backoff{intervals: ret}
}

type backoff struct {
	index     int
	len       int
	intervals []time.Duration
	mux       sync.Mutex
}

func (b *backoff) Next() {
	b.mux.Lock()
	defer b.mux.Unlock()

	// memoize intervals len
	if b.len == 0 {
		b.len = len(b.intervals)
	}

	if b.index+1 >= b.len {
		return
	}

	b.index++
}

func (b *backoff) Current() time.Duration {
	return b.intervals[b.index]
}

func (b *backoff) Reset() {
	b.mux.Lock()
	defer b.mux.Unlock()

	b.index = 0
}
