package main

import (
	"log"
	"os"
)

var DEBUG = 0
var ERROR = 10

// default level is ERROR
var logLevel = ERROR

var LOG = log.New(os.Stdout, "[DEBUG] ", log.Ltime)
var ERR = log.New(os.Stdout, "[ERROR] ", log.Ltime)

func SetLogLevel(level int) {
	logLevel = level
}

func Debugln(args ...interface{}) {
	if logLevel == DEBUG {
		LOG.Println(args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if logLevel == DEBUG {
		LOG.Printf(format, args...)
	}
}

func Fatalln(args ...interface{}) {
	LOG.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	LOG.Fatalf(format, args...)
}

func Errorln(args ...interface{}) {
	ERR.Println(args...)
}

func Errorf(format string, args ...interface{}) {
	ERR.Printf(format, args...)
}
