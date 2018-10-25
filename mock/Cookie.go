package mock

import "github.com/cookiejars/cookiejar"

var _ cookiejar.Cookie = &Cookie{}

type Cookie struct {
	ContentFn      func() ([]byte, error)
	ContentInvoked bool

	DeleteFn      func() error
	DeleteInvoked bool
}

func (c *Cookie) Content() ([]byte, error) {
	c.ContentInvoked = true
	return c.ContentFn()
}

func (c *Cookie) Delete() error {
	c.DeleteInvoked = true
	return c.DeleteFn()
}
