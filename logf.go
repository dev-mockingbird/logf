package logf

import (
	"fmt"
	"log"
	"os"
	"strings"
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

type Prefixable interface {
	Prefix(string) Logger
}

type Logger interface {
	Prefixable
	Logfer
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
	prefixes   []string
	logLevel   Level
}

type Option func(*logger)

func LogLevel(logLevel Level) Option {
	return func(opt *logger) {
		opt.logLevel = logLevel
	}
}

func Prefix(prefix string) Option {
	return func(opt *logger) {
		opt.prefixes = append(opt.prefixes, prefix)
	}
}

func Underlying(underlying *log.Logger) Option {
	return func(opt *logger) {
		opt.underlying = underlying
	}
}

func New(opts ...Option) Logger {
	logger := &logger{logLevel: Info}
	for _, apply := range opts {
		apply(logger)
	}
	logger.underlying = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	return logger
}

func (l logger) Logf(level Level, format string, v ...any) {
	if level >= l.logLevel {
		l.underlying.Output(2, fmt.Sprintf("[%s] %s%s", LevelString(level), strings.Join(l.prefixes, ""), fmt.Sprintf(format, v...)))
	}
}

func (l *logger) Prefix(prefix string) Logger {
	ret := *l
	ret.prefixes = append(l.prefixes, prefix)
	return &ret
}

func LevelString(level Level) string {
	if s, ok := levelString[level]; ok {
		return s
	}
	return "UNKOWN"
}
