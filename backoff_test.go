package cookiejar_test

import (
	"testing"
	"time"

	"github.com/cookiejars/cookiejar"
)

func TestConstantBackoff(t *testing.T) {
	t.Parallel()

	b := cookiejar.ConstantBackoff(time.Second)
	iterations := []iteration{
		{fn: func() {}, expectedCurrent: time.Duration(0)},
		{fn: b.Next, expectedCurrent: time.Second},
		{fn: b.Next, expectedCurrent: time.Second},
		{fn: b.Next, expectedCurrent: time.Second},
		{fn: b.Reset, expectedCurrent: time.Duration(0)},
	}

	testIterations(t, b.Current, iterations...)
}

func TestExponentialBackoff(t *testing.T) {
	t.Parallel()

	b := cookiejar.ExponentialBackoff(2, time.Second)
	iterations := []iteration{
		{fn: func() {}, expectedCurrent: time.Duration(0)},
		{fn: b.Next, expectedCurrent: time.Second},
		{fn: b.Next, expectedCurrent: time.Second * 2},
		{fn: b.Next, expectedCurrent: time.Second * 2},
		{fn: b.Next, expectedCurrent: time.Second * 2},
		{fn: b.Reset, expectedCurrent: time.Duration(0)},
	}

	testIterations(t, b.Current, iterations...)
}

type iteration struct {
	fn              func()
	expectedCurrent time.Duration
}

func testIterations(t *testing.T, getCurrent func() time.Duration, iterations ...iteration) {
	for _, it := range iterations {
		it.fn()
		assertEqual(t, getCurrent(), it.expectedCurrent)
	}
}
