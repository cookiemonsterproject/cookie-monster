package main

import (
	"fmt"
	"path/filepath"
	"syscall"

	"github.com/cookiejars/cookiejar"
)

func digestFn() cookiejar.DigestFn {
	return func(cookie cookiejar.Cookie) error {
		fmt.Println(cookie.Content().(string))
		return nil
	}
}

func main() {
	path, err := filepath.Abs("jar.so")
	if err != nil {
		panic(err)
	}

	dig, err := cookiejar.NewDigesterWithPlugin(
		path,
		cookiejar.SetInfoLog(infoLogger{}),
		cookiejar.SetErrorLog(errorLogger{}),
		cookiejar.SetStopSignals(syscall.SIGINT),
	)
	if err != nil {
		panic(err)
	}

	err = dig.Start(digestFn())
	if err != nil {
		panic(err)
	}
}
