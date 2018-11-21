package main

import "fmt"

type infoLogger struct{}

func (infoLogger) Printf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

type errorLogger struct{}

func (errorLogger) Printf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
