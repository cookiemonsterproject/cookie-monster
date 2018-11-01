package mock

import "github.com/cookiejars/cookiejar"

var _ cookiejar.Jar = &Jar{}

type Jar struct {
	RetrieveFn      func() ([]cookiejar.Cookie, error)
	RetrieveInvoked bool

	RetireFn      func(cookie cookiejar.Cookie) error
	RetireInvoked bool
}

func (j *Jar) Retrieve() ([]cookiejar.Cookie, error) {
	j.RetrieveInvoked = true
	return j.RetrieveFn()
}

func (j *Jar) Retire(cookie cookiejar.Cookie) error {
	j.RetireInvoked = true
	return j.RetireFn(cookie)
}
