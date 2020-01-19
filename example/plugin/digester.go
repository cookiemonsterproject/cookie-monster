package main

import (
	"fmt"
	"path/filepath"
	"syscall"

	"github.com/cookiejars/cookiemonster"
)

func digestFn() cookiemonster.DigestFn {
	return func(cookie cookiemonster.Cookie) error {
		fmt.Println(cookie.Content().(string))
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

	err = dig.Start(digestFn())
	if err != nil {
		panic(err)
	}
}
