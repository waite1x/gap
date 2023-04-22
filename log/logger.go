package log

import (
	"fmt"
	"runtime/debug"
)

type Level int8

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
	NoLevel
	Disabled
	TraceLevel Level = -1
)

type Options struct {
	Level Level
}

var Opts = &Options{
	Level: InfoLevel,
}

var _logger Logger = &DefaultLogger{
	opts: Opts,
}

type Logger interface {
	Log(message string, level Level)
}

func SetLogger(logger Logger) {
	_logger = logger
}

func GetLogger() Logger {
	return _logger
}

func Info(msg string) {
	GetLogger().Log(msg, InfoLevel)
}

func Infof(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	GetLogger().Log(msg, InfoLevel)
}

func Debug(msg string) {
	GetLogger().Log(msg, DebugLevel)
}

func Debugf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	GetLogger().Log(msg, DebugLevel)
}

func Warn(msg string) {
	GetLogger().Log(msg, WarnLevel)
}

func Warnf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	GetLogger().Log(msg, WarnLevel)
}

func Error(msg string, err error) {
	stack := string(debug.Stack())
	msg = fmt.Sprintf("%s, error: %s\n%s", msg, err, stack)
	GetLogger().Log(msg, ErrorLevel)
}
