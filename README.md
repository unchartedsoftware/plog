# Plog

> A pretty logger for Go

## Dependencies

Requires the [Go](https://golang.org/) programming language binaries with the `GOPATH` environment variable specified.

## Installation

```bash
go get github.com/unchartedsoftware/plog
```

## Example

This minimalistic application shows how to log at different levels and set a log level filter.

```go
package main

import (
    log "github.com/unchartedsoftware/plog"
)

func main() {   
    log.Debug("This is a debug level log")
    log.Info("This is an info level log")
    log.Warn("This is a warn level log")
    log.Error("This is an error level log")

    // only log warnings and errors
    log.SetLevel(log.WarnLevel)

    log.Debug("This is an info level log, I will be ignored")
    log.Info("This is a debug level log, I too shall be ignored")
    log.Warn("This is a warn level log, you will see me")
    log.Error("This is an error level log, you will see me")
}
```
