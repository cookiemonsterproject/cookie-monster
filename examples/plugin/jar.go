package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cookiemonsterproject/cookie-monster"
)

// Jar is the plugin's exported symbol
var Jar jar

type jar struct{}

func (j jar) Retrieve() ([]cookiemonster.Cookie, error) {
	// try to simulate a real system
	if time.Now().Nanosecond()%2 == 0 {
		return nil, nil
	}

	return getCookies(), nil
}

func (jar) Retire(cookiemonster.Cookie) error {
	return nil
}

func getCookies() []cookiemonster.Cookie {
	cookies := make([]cookiemonster.Cookie, rand.Intn(10))

	for i, _ := range cookies {
		now := time.Now()
		cookie := c{
			id:      fmt.Sprintf("id-%d", now.Unix()),
			content: now.String(),
		}

		cookies[i] = cookie
	}

	return cookies
}

type c struct {
	id      string
	content string
}

func (c c) ID() string {
	return c.id
}

func (c c) Content() interface{} {
	return c.content
}

func (c c) Metadata() map[string]string {
	return nil
}
