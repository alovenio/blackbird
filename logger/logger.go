package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	Debug = iota
	Info
	Warn
	Error
)

var LogLevel = Info

var debugLogger = log.New(os.Stdout, "[DEBUG] ", 0xF)
var infoLogger = log.New(os.Stdout, "\u001B[36m[INFO] \u001B[0m", 0xF)
var warnLogger = log.New(os.Stdout, "\u001B[33m[WARN] \u001B[0m", 0xF)
var errorLogger = log.New(os.Stdout, "\u001b[31m[ERROR] \u001b[0m", 0xF)

func ParseLogLevel(level string) (int, error) {
	var l = Info
	var e error
	switch strings.ToLower(level) {
	case "info":
		l = Info
	case "debug":
		l = Debug
	case "error":
		l = Error
	case "warn":
		l = Warn
	default:
		e = fmt.Errorf("unknown log level %q", level)
	}
	return l, e
}

func LogDebugF(fmt string, args ...any) {
	if LogLevel <= Debug {
		debugLogger.Printf(fmt, args...)
	}
}

func LogInfoF(fmt string, args ...any) {
	if LogLevel <= Info {
		infoLogger.Printf(fmt, args...)
	}
}

func LogWarnF(fmt string, args ...any) {
	if LogLevel <= Warn {
		warnLogger.Printf(fmt, args...)
	}
}

func LogErrorF(fmt string, args ...any) {
	if LogLevel <= Error {
		errorLogger.Printf(fmt, args...)
	}
}

func LogFatalF(err error) {
	errorLogger.Println(err)
	os.Exit(1)
}
