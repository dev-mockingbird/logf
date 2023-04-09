package logf

import "testing"

func TestLogger(t *testing.T) {
	logger := New()
	logger.Logf(Trace, "hello world: %d", 1)
	logger.Logf(Debug, "hello world: %s", "yang,zhong")
	logger.Logf(Info, "hello world: %v", struct{ Name string }{Name: "yang,zhong"})
	logger.Logf(Warn, "hello world: %#v", struct{ Name string }{Name: "yang,zhong"})
	logger.Logf(Error, "hello world")
	logger.Logf(Fatal, "hello world")
}
