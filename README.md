# Plog

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](http://godoc.org/github.com/unchartedsoftware/plog)
[![Build Status](https://travis-ci.org/unchartedsoftware/plog.svg?branch=master)](https://travis-ci.org/unchartedsoftware/plog)
[![Go Report Card](https://goreportcard.com/badge/github.com/unchartedsoftware/plog)](https://goreportcard.com/report/github.com/unchartedsoftware/plog)

> A pretty logger for Go

## Example

```go
package main

import (
	"github.com/unchartedsoftware/plog"
)

func main() {
	// use global log methods
	log.Debug("This is a debug level log")
	log.Info("This", "is", "an", "info", "level", "log")
	log.Warnf("This is a %s level log", "warn")
	log.Error("This is an error level log")

	// only log warnings and errors
	log.SetLevel(log.WarnLevel)

	log.Debug("This is a debug level log, I will be ignored")
	log.Info("This is an info level log, I will be ignored")
	log.Warn("This is a warn level log, you will see me")
	log.Error("This is an error level log, you will see me")

	// create logger instance
	logger := log.NewLogger()
	logger.SetLevel(log.ErrorLevel)
	logger.Debug("This is a debug level log, I will be ignored")
	logger.Info("This is an info level log, I will be ignored")
	logger.Warn("This is a warn level log, I will be ignored")
	logger.Error("This is an error level log, you will see me")
}
```
