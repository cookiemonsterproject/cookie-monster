package cookiejar_test

import (
	"testing"
	"time"

	"github.com/cookiejars/cookiejar"
)

func TestConstantBackoff(t *testing.T) {
	t.Parallel()

	b := cookiejar.ConstantBackoff(time.Second)
	assertEqual(t, b.Current(), time.Duration(0))

	b.Next()
	assertEqual(t, b.Current(), time.Second)

	b.Next()
	assertEqual(t, b.Current(), time.Second)

	b.Reset()
	assertEqual(t, b.Current(), time.Duration(0))
}

func TestExponentialBackoff(t *testing.T) {
	t.Parallel()

	b := cookiejar.ExponentialBackoff(2, time.Second)
	assertEqual(t, b.Current(), time.Duration(0))

	b.Next()
	assertEqual(t, b.Current(), time.Second)

	b.Next()
	assertEqual(t, b.Current(), time.Second*2)

	b.Next()
	assertEqual(t, b.Current(), time.Second*2)

	b.Reset()
	assertEqual(t, b.Current(), time.Duration(0))
}
