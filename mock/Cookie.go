package mock

import "github.com/cookiejars/cookiejar"

var _ cookiejar.Cookie = &Cookie{}

type Cookie struct {
	ContentFn      func() (interface{}, error)
	ContentInvoked bool
}

func (c *Cookie) Content() (interface{}, error) {
	c.ContentInvoked = true
	return c.ContentFn()
}
