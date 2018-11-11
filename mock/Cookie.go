package mock

import "github.com/cookiejars/cookiejar"

var _ cookiejar.Cookie = &Cookie{}

type Cookie struct {
	IDFn      func() string
	IDInvoked bool

	ContentFn      func() interface{}
	ContentInvoked bool
}

func (c *Cookie) ID() string {
	c.IDInvoked = true
	return c.IDFn()
}

func (c *Cookie) Content() interface{} {
	c.ContentInvoked = true
	return c.ContentFn()
}
