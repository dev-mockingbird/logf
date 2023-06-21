package logf

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"github.com/fatih/color"
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
	underlying   *log.Logger
	disableColor bool
	prefixes     []string
	logLevel     Level
}

type Option func(*logger)

func LogLevel(logLevel Level) Option {
	return func(opt *logger) {
		opt.logLevel = logLevel
	}
}

func DesiableColor() Option {
	return func(l *logger) {
		l.disableColor = true
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
	if level < l.logLevel {
		return
	}
	ls := LevelString(level)
	msg := fmt.Sprintf("\t[%s]\t%s%s", ls, strings.Join(l.prefixes, ""), fmt.Sprintf(format, v...))
	l.underlying.Output(2, l.colorMsg(level, msg))
	if level >= Warn {
		stack := debug.Stack()
		stacks := bytes.Split(stack, []byte{'\n'})[5:]
		l.underlying.Output(2, l.colorMsg(level, string(bytes.Join(stacks, []byte{'\n'}))))
	}
}

func (l logger) colorMsg(level Level, msg string) string {
	if l.disableColor {
		return msg
	}
	switch level {
	case Trace:
		return color.WhiteString(msg)
	case Debug:
		return color.CyanString(msg)
	case Info:
		return color.GreenString(msg)
	case Warn:
		return color.YellowString(msg)
	case Error:
		return color.RedString(msg)
	case Fatal:
		return color.BlueString(msg)
	}
	return msg
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
