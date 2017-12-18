package main

import (
	"github.com/op/go-logging"
)

const (
	log_name = "hass"
)

var (
	format = logging.MustStringFormatter(
		`%{color}%{level:.5s}%{color:reset} %{message}`,
	)
	log = logging.MustGetLogger(log_name)
)

func init() {
	logging.SetFormatter(format)
	logging.SetLevel(logging.INFO, log_name)
}

func EnableDebug() {
	logging.SetLevel(logging.DEBUG, log_name)
}

func Debugln(args ...interface{}) {
	log.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Infoln(args ...interface{}) {
	log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Fatalln(args ...interface{}) {
	log.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Errorln(args ...interface{}) {
	log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}
