package main

import (
	"fmt"
	"time"

	"github.com/cookiejars/cookiemonster"
)

// Jar is the plugin's exported symbol
var Jar jar

type jar struct{}

func (j jar) Retrieve() ([]cookiemonster.Cookie, error) {
	// simulate a real system
	if time.Now().Second()%2 == 0 {
		return nil, nil
	}

	return getCookies(), nil
}

func (jar) Retire(cookiemonster.Cookie) error {
	return nil
}

func getCookies() []cookiemonster.Cookie {
	now := time.Now()

	cookie := c{
		id:      fmt.Sprintf("id-%d", now.Unix()),
		content: now.String(),
	}

	return []cookiemonster.Cookie{cookie}
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
