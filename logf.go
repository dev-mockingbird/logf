package logf

import (
	"fmt"
	"log"
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
	underling *log.Logger
}

func New(loggers ...*log.Logger) Logfer {
	if len(loggers) > 0 {
		return logger{underling: loggers[0]}
	}
	return logger{underling: log.Default()}
}

func (l logger) Logf(level Level, format string, v ...any) {
	l.underling.Printf("%s  %s", LevelString(level), fmt.Sprintf(format, v...))
}

func LevelString(level Level) string {
	if s, ok := levelString[level]; ok {
		return s
	}
	return "UNKOWN"
}
