package logf

import (
	"fmt"
	"log"
	"os"
)

type Level uint8

const (
	Trace Level = iota
	Debug
	Info
	Warn
	Error
	Fatal
)

var levelString map[Level]string

func init() {
	levelString = map[Level]string{
		Trace: "TRACE",
		Debug: "DEBUG",
		Info:  "INFO",
		Warn:  "WARN",
		Error: "ERROR",
		Fatal: "FATAL",
	}
}

type Logfer interface {
	Logf(level Level, format string, v ...interface{})
}

type Logf func(level Level, format string, v ...interface{})

func (logf Logf) Logf(level Level, format string, v ...interface{}) {
	logf(level, format, v...)
}

type logger struct {
	underlying *log.Logger
	logLevel   Level
}

type Option func(*logger)

func LogLevel(logLevel Level) Option {
	return func(opt *logger) {
		opt.logLevel = logLevel
	}
}

func Underlying(underlying *log.Logger) Option {
	return func(opt *logger) {
		opt.underlying = underlying
	}
}

func New(opts ...Option) Logfer {
	logger := &logger{underlying: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile), logLevel: Trace}
	for _, apply := range opts {
		apply(logger)
	}
	return logger
}

func (l logger) Logf(level Level, format string, v ...any) {
	if level >= l.logLevel {
		l.underlying.Printf("[%s] %s", LevelString(level), fmt.Sprintf(format, v...))
	}
}

func LevelString(level Level) string {
	if s, ok := levelString[level]; ok {
		return s
	}
	return "UNKOWN"
}
