package main

import (
	"fmt"
	"time"

	"github.com/cookiejars/cookiejar"
)

// Jar is the plugin's exported symbol
var Jar jar

type jar struct{}

func (j jar) Retrieve() ([]cookiejar.Cookie, error) {
	// simulate a real system
	if time.Now().Second()%2 == 0 {
		return nil, nil
	}

	return getCookies(), nil
}

func (jar) Retire(cookiejar.Cookie) error {
	return nil
}

func getCookies() []cookiejar.Cookie {
	now := time.Now()

	cookie := c{
		id:      fmt.Sprintf("id-%d", now.Unix()),
		content: now.String(),
	}

	return []cookiejar.Cookie{cookie}
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
