package logf

import "testing"

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

func TestLogger_Prefix(t *testing.T) {
	logger := New(LogLevel(Info), Prefix("prefix xxxx: "))
	logger.Logf(Trace, "hello world: %d", 1)
	logger.Logf(Debug, "hello world: %s", "yang,zhong")
	logger.Logf(Info, "hello world: %v", struct{ Name string }{Name: "yang,zhong"})
	logger.Logf(Warn, "hello world: %#v", struct{ Name string }{Name: "yang,zhong"})
	logger.Logf(Error, "hello world")
	logger.Logf(Fatal, "hello world")
}
