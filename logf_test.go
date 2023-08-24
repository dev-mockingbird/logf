package logf

import (
	"bytes"
	"strings"
	"testing"
)

func _log(level Level, format string, args ...any) {
	logger := New(LogLevel(Trace), WithCallerDepth(CallerDepth+1))
	logger.Logf(level, format, args...)
}

func TestLogger_CallerDepth(t *testing.T) {
	_log(Info, "hello world")
	_log(Warn, "hello world")
}

func TestLogger(t *testing.T) {
	logger := New(LogLevel(Trace))
	logger.Logf(Trace, "hello world: %d", 1)
	logger.Logf(Debug, "hello world: %s", "yang,zhong")
	logger.Logf(Info, "hello world: %v", struct{ Name string }{Name: "yang,zhong"})
	logger.Logf(Warn, "hello world: %#v", struct{ Name string }{Name: "yang,zhong"})
	logger.Logf(Error, "hello world")
	logger.Logf(Fatal, "hello world")
}

func TestLogger_LogLevel(t *testing.T) {
	logger := New(LogLevel(Info))
	logger.Logf(Trace, "hello world: %d", 1)
	logger.Logf(Debug, "hello world: %s", "yang,zhong")
	logger.Logf(Info, "hello world: %v", struct{ Name string }{Name: "yang,zhong"})
	logger.Logf(Warn, "hello world: %#v", struct{ Name string }{Name: "yang,zhong"})
	logger.Logf(Error, "hello world")
	logger.Logf(Fatal, "hello world")
}

func TestLogger_LogToCustomWriter(t *testing.T) {
	var buf bytes.Buffer
	logger := New(LogLevel(Info), Writer(&buf))
	logger.Logf(Info, "hello world")
	result := buf.String()
	if !strings.Contains(result, "hello world") {
		t.Fail()
	}
}

func TestLogger_Prefix(t *testing.T) {
	logger := New(LogLevel(Info), Prefix("prefix xxxx: "))
	logger.Logf(Trace, "hello world: %d", 1)
	logger.Logf(Debug, "hello world: %s", "yang,zhong")
	logger.Logf(Info, "hello world: %v", struct{ Name string }{Name: "yang,zhong"})
	logger.Logf(Warn, "hello world: %#v", struct{ Name string }{Name: "yang,zhong"})
	logger.Logf(Error, "hello world")
	logger.Logf(Fatal, "hello world")
}
