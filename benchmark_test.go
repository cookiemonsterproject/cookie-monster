package cookiejar

import (
	"fmt"
	"math"
	"testing"
)

func BenchmarkDigester(b *testing.B) {
	for k := 0.; k <= 15; k++ {
		n := int(math.Pow(2, k))

		digester := getInitializedDigester(n)

		b.Run(fmt.Sprintf("%d-workers", n), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				digester.workChan <- []Cookie{}
			}
		})

		digester.Stop()
	}
}

func getInitializedDigester(workers int) *digester {
	digestFn := func(cookie Cookie) error {
		return nil
	}

	d := NewDigester(workers, nil, nil).(*digester)
	d.startWorkers(digestFn)

	return d
}
