package log

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
)

const (
	// DebugLevel logging is for development level logging.
	DebugLevel = 1
	// InfoLevel logging is for high granularity development logging events.
	InfoLevel = 2
	// WarnLevel logging is for unexpected and recoverable events.
	WarnLevel = 3
	// ErrorLevel logging is for unexpected and unrecoverable fatal events.
	ErrorLevel = 4
)

var (
	log = getLogger()
	showRoutineID = false
)

func getLogger() *logrus.Logger {
	log := logrus.New()
	log.Formatter = new(PrettyFormatter)
	log.Level = logrus.DebugLevel
	return log
}

func getGoroutineID() (uint64, error) {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	return strconv.ParseUint(string(b), 10, 64)
}

func retrieveCallInfo() string {
	pc, file, line, _ := runtime.Caller(2)
	// get full path to file
	fullpath := runtime.FuncForPC(pc).Name()
	// strip out vendor path if it exists
	parts := strings.Split(fullpath, "/vendor/")
	pckg := parts[len(parts) - 1]
	// get package name
	parts = strings.Split(pckg, "/")
	lastIndex := len(parts) - 1
	index := 3 // domain/company/root
	if index > lastIndex {
		index = lastIndex
	}
	// remove function
	parts[lastIndex] = strings.Split(parts[lastIndex], ".")[0]
	packageName := strings.Join(parts[index:], "/")
	// get file name
	_, fileName := path.Split(file)
	// determine whether or not to show goroutine id
	if (showRoutineID) {
		gid, err := getGoroutineID()
		if err == nil {
			return fmt.Sprint(packageName, "/", fileName, ":", line, ", gid:", gid)
		}
	}
	return fmt.Sprint(packageName, "/", fileName, ":", line)
}

// Infof logging is for high granularity development logging events.
func Infof(format string, args ...interface{}) {
	log.WithFields(logrus.Fields{
		"fileinfo": retrieveCallInfo(),
	}).Infof(format, args...)
}

// Debugf logging is for development level logging events.
func Debugf(format string, args ...interface{}) {
	log.WithFields(logrus.Fields{
		"fileinfo": retrieveCallInfo(),
	}).Debugf(format, args...)
}

// Warnf logging is for unexpected and recoverable events.
func Warnf(format string, args ...interface{}) {
	log.WithFields(logrus.Fields{
		"fileinfo": retrieveCallInfo(),
	}).Warnf(format, args...)
}

// Errorf level is for unexpected and unrecoverable fatal events.
func Errorf(format string, args ...interface{}) {
	log.WithFields(logrus.Fields{
		"fileinfo": retrieveCallInfo(),
	}).Errorf(format, args...)
}

// Info logging is for high granularity development logging events.
func Info(args ...interface{}) {
	log.WithFields(logrus.Fields{
		"fileinfo": retrieveCallInfo(),
	}).Info(args...)
}

// Debug logging is for development level logging events.
func Debug(args ...interface{}) {
	log.WithFields(logrus.Fields{
		"fileinfo": retrieveCallInfo(),
	}).Debug(args...)
}

// Warn logging is for unexpected and recoverable events.
func Warn(args ...interface{}) {
	log.WithFields(logrus.Fields{
		"fileinfo": retrieveCallInfo(),
	}).Warn(args...)
}

// Error level is for unexpected and unrecoverable fatal events.
func Error(args ...interface{}) {
	log.WithFields(logrus.Fields{
		"fileinfo": retrieveCallInfo(),
	}).Error(args...)
}

// SetLevel sets the current logging output level.
func SetLevel(level int) {
	switch level {
	case DebugLevel:
		log.Level = logrus.DebugLevel
	case InfoLevel:
		log.Level = logrus.InfoLevel
	case WarnLevel:
		log.Level = logrus.WarnLevel
	case ErrorLevel:
		log.Level = logrus.ErrorLevel
	}
}

func ShowGoRoutineID() {
	showRoutineID = true
}

func HideGoRoutineID() {
	showRoutineID = false
}
