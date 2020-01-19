package mock

import "github.com/cookiejars/cookiemonster"

var _ cookiemonster.Cookie = &Cookie{}

type Cookie struct {
	IDFn      func() string
	IDInvoked bool

	ContentFn      func() interface{}
	ContentInvoked bool

	MetadataFn      func() map[string]string
	MetadataInvoked bool
}

func (c *Cookie) ID() string {
	c.IDInvoked = true
	return c.IDFn()
}

func (c *Cookie) Content() interface{} {
	c.ContentInvoked = true
	return c.ContentFn()
}

func (c *Cookie) Metadata() map[string]string {
	c.MetadataInvoked = true
	return c.MetadataFn()
}
