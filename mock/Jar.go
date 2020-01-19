package mock

import "github.com/cookiemonsterproject/cookie-monster"

var _ cookiemonster.Jar = &Jar{}

type Jar struct {
	RetrieveFn      func() ([]cookiemonster.Cookie, error)
	RetrieveInvoked bool

	RetireFn      func(cookie cookiemonster.Cookie) error
	RetireInvoked bool
}

func (j *Jar) Retrieve() ([]cookiemonster.Cookie, error) {
	j.RetrieveInvoked = true
	return j.RetrieveFn()
}

func (j *Jar) Retire(cookie cookiemonster.Cookie) error {
	j.RetireInvoked = true
	return j.RetireFn(cookie)
}
