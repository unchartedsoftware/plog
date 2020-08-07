package log

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var (
	logger = NewLogger()
	output = os.Stdout
)

const (
	// DebugLevel logging is for high granularity development level logging.
	DebugLevel Level = 1
	// InfoLevel logging is for logging events.
	InfoLevel Level = 2
	// WarnLevel logging is for unexpected and recoverable events.
	WarnLevel Level = 3
	// ErrorLevel logging is for unexpected and unrecoverable fatal events.
	ErrorLevel Level = 4
	// the depth required to provide the correct caller info
	defaultDepth = 3
)

// Level represents the logging level type enumeration.
type Level int

// Logger represents a basic logging struct.
type Logger struct {
	showRoutineID bool
	loggingLevel  Level
	depth         int
	mu            *sync.Mutex
}

func init() {
	logger = NewLogger()
	logger.depth = defaultDepth + 1 // increment this for global logger
}

// NewLogger instantiates and returns a new Logger struct.
func NewLogger() *Logger {
	return &Logger{
		loggingLevel: DebugLevel,
		mu:           &sync.Mutex{},
		depth:        defaultDepth,
	}
}

// Debugf is for debug level logging events.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.writeOutputf(l.depth, DebugLevel, format, args...)
}

// Infof is for high granularity logging events.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.writeOutputf(l.depth, InfoLevel, format, args...)
}

// Warnf is for unexpected and recoverable events.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.writeOutputf(l.depth, WarnLevel, format, args...)
}

// Errorf is for unexpected and unrecoverable fatal events.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.writeOutputf(l.depth, ErrorLevel, format, args...)
}

// Debug is for debug level logging events.
func (l *Logger) Debug(args ...interface{}) {
	l.writeOutput(l.depth, DebugLevel, args...)
}

// Info is for high granularity logging events.
func (l *Logger) Info(args ...interface{}) {
	l.writeOutput(l.depth, InfoLevel, args...)
}

// Warn is for unexpected and recoverable events.
func (l *Logger) Warn(args ...interface{}) {
	l.writeOutput(l.depth, WarnLevel, args...)
}

// Error is for unexpected and unrecoverable fatal events.
func (l *Logger) Error(args ...interface{}) {
	l.writeOutput(l.depth, ErrorLevel, args...)
}

// SetLevel sets the current logging output level.
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	// only set if valid argument
	switch level {
	case DebugLevel:
		l.loggingLevel = level
	case InfoLevel:
		l.loggingLevel = level
	case WarnLevel:
		l.loggingLevel = level
	case ErrorLevel:
		l.loggingLevel = level
	}
}

// ShowGoRoutineID enables appending the goroutine ID to the log output.
func (l *Logger) ShowGoRoutineID() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.showRoutineID = true
}

// HideGoRoutineID disables appending the goroutine ID to the log output.
func (l *Logger) HideGoRoutineID() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.showRoutineID = false
}

func (l *Logger) getGoroutineID() (uint64, error) {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	return strconv.ParseUint(string(b), 10, 64)
}

func (l *Logger) retrieveCallInfo(depth int) string {
	pc, file, line, _ := runtime.Caller(depth)
	// get full path to file
	fullpath := runtime.FuncForPC(pc).Name()
	// strip out vendor path if it exists
	parts := strings.Split(fullpath, "/vendor/")
	pckg := parts[len(parts)-1]
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
	if l.showRoutineID {
		gid, err := l.getGoroutineID()
		if err == nil {
			return fmt.Sprint(packageName, "/", fileName, ":", line, ", gid:", gid)
		}
	}
	return fmt.Sprint(packageName, "/", fileName, ":", line)
}

func (l *Logger) sprint(args ...interface{}) string {
	// ensure that spaces are put between ALL operands
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}

func (l *Logger) writeOutputf(depth int, level Level, format string, args ...interface{}) {
	if level < l.loggingLevel {
		return
	}
	l.mu.Lock()
	writer := bufio.NewWriter(output)
	msg := fmt.Sprintf(format, args...)
	defer l.mu.Unlock()  // then unlock
	defer writer.Flush() // flush first
	writer.Write(formatLog(level, msg, l.retrieveCallInfo(depth)))
}

func (l *Logger) writeOutput(depth int, level Level, args ...interface{}) {
	if level < l.loggingLevel {
		return
	}
	l.mu.Lock()
	writer := bufio.NewWriter(output)
	defer l.mu.Unlock()  // then unlock
	defer writer.Flush() // flush first
	writer.Write(formatLog(level, l.sprint(args...), l.retrieveCallInfo(depth)))
}

// Debugf is for debug level logging events.
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Infof is for high granularity logging events.
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Warnf is for unexpected and recoverable events.
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Errorf is for unexpected and unrecoverable fatal events.
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Debug is for debug level logging events.
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Info is for high granularity logging events.
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Warn is for unexpected and recoverable events.
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Error is for unexpected and unrecoverable fatal events.
func Error(args ...interface{}) {
	logger.Error(args...)
}

// SetLevel sets the current logging output level.
func SetLevel(level Level) {
	logger.SetLevel(level)
}

// ShowGoRoutineID enables appending the goroutine ID to the log output.
func ShowGoRoutineID() {
	logger.ShowGoRoutineID()
}

// HideGoRoutineID disables appending the goroutine ID to the log output.
func HideGoRoutineID() {
	logger.HideGoRoutineID()
}
