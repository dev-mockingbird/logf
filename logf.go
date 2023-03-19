package logf

import "go-micro.dev/v4/logger"

type Logfer interface {
	Logf(level logger.Level, format string, v ...interface{})
}

type Logf func(level logger.Level, format string, v ...interface{})

func (logf Logf) Logf(level logger.Level, format string, v ...interface{}) {
	logf(level, format, v...)
}

func New() Logfer {
	return logger.NewLogger()
}
