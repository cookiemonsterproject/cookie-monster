package cookiejar

type DigestFn func(cookie Cookie) error

type Digester interface {
	Start(fn DigestFn) error
}

type digester struct {
	workers int

	backoff Backoff
}

func NewDigester(workers int, backoff Backoff) Digester {
	return digester{
		workers: workers,
		backoff: backoff,
	}
}

func (digester) Start(fn DigestFn) error {
	return nil
}

type Cookie interface {
	Content() ([]byte, error)
}

type Jar interface {
	Fetch() (Cookie, error)
}
