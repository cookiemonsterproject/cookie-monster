package mock

import "github.com/cookiejars/cookiejar"

var _ cookiejar.Jar = &Jar{}

type Jar struct {
	RetrieveFn      func() ([]cookiejar.Cookie, error)
	RetrieveInvoked bool
}

func (j *Jar) Retrieve() ([]cookiejar.Cookie, error) {
	j.RetrieveInvoked = true
	return j.RetrieveFn()
}
