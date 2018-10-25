package mock

import "github.com/cookiejars/cookiejar"

var _ cookiejar.Jar = &Jar{}

type Jar struct {
	FetchFn      func() ([]cookiejar.Cookie, error)
	FetchInvoked bool
}

func (j *Jar) Fetch() ([]cookiejar.Cookie, error) {
	j.FetchInvoked = true
	return j.FetchFn()
}
