package main

import (
	"fmt"
	"path/filepath"
	"syscall"

	"github.com/cookiemonsterproject/cookie-monster"
)

func digestFn() cookiemonster.DigestFn {
	return func(cookie cookiemonster.Cookie) error {
		fmt.Printf("Cookie's content: %s\n", cookie.Content().(string))
		return nil
	}
}

func main() {
	path, err := filepath.Abs("jar.so")
	if err != nil {
		panic(err)
	}

	dig, err := cookiemonster.NewDigesterWithPlugin(
		path,
		cookiemonster.SetInfoLog(infoLogger{}),
		cookiemonster.SetErrorLog(errorLogger{}),
		cookiemonster.SetStopSignals(syscall.SIGINT),
	)
	if err != nil {
		panic(err)
	}

	if err = dig.Start(digestFn()); err != nil {
		panic(err)
	}
}
