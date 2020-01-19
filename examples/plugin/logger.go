package main

import (
	"log"
)

type infoLogger struct{}

func (infoLogger) Printf(format string, args ...interface{}) {
	log.Printf(format+"\n", args...)
}

type errorLogger struct{}

func (errorLogger) Printf(format string, args ...interface{}) {
	log.Printf(format+"\n", args...)
}
