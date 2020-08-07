package log

import (
	"bytes"
	"fmt"
	"regexp"
	"runtime"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/mgutz/ansi"
)

var (
	isTerminal = isatty.IsTerminal(output.Fd())
	isColored  = isTerminal && (runtime.GOOS != "windows")
	re         = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")
)

func formatLog(level Level, msg string, fileinfo string) []byte {
	b := &bytes.Buffer{}
	if isColored {
		printColored(b, level, msg, fileinfo)
	} else {
		printUncolored(b, level, msg, fileinfo)
	}
	b.WriteByte('\n')
	return b.Bytes()
}

func getLevelString(level Level) string {
	switch level {
	case InfoLevel:
		return " INFO "
	case WarnLevel:
		return " WARN "
	case ErrorLevel:
		return " ERROR "
	default:
		return " DEBUG "
	}
}

func getLevelColor(level Level) string {
	switch level {
	case InfoLevel:
		return ansi.Reset
	case WarnLevel:
		return ansi.Yellow
	case ErrorLevel:
		return ansi.Red
	default:
		return ansi.Blue
	}
}

func printColored(b *bytes.Buffer, level Level, msg string, fileinfo string) {
	levelText := getLevelString(level)
	levelColor := getLevelColor(level)
	// write log message to buffer
	fmt.Fprintf(b, "%s[ %s ]%s %s[%s]%s %s %s(%s)%s",
		ansi.LightBlack,
		time.Now().Format(time.Stamp),
		ansi.Reset,
		levelColor,
		levelText,
		ansi.Reset,
		msg,
		ansi.Cyan,
		fileinfo,
		ansi.Reset)
}

func stripColor(str string) string {
	return re.ReplaceAllString(str, "")
}

func printUncolored(b *bytes.Buffer, level Level, msg string, fileinfo string) {
	levelText := getLevelString(level)
	// write log message to buffer
	fmt.Fprintf(b, "[ %s ] [%s] %s (%s)",
		time.Now().Format(time.Stamp),
		levelText,
		stripColor(msg),
		fileinfo)
}
