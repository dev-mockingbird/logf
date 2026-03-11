package logf

import (
	"bytes"
	"strings"
	"testing"
)

func _log(level Level, format string, args ...any) {
	logger := New(LogLevel(Trace), Caller(CallerDepth+1))
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

func TestLogger_MergeStackInOneLine(t *testing.T) {
	var buf bytes.Buffer
	// 使用高于 Warn 的级别以触发 Stack 收集
	logger := New(LogLevel(Warn), Writer(&buf))
	logger.Logf(Warn, "main message")

	result := buf.String()
	lines := strings.Split(strings.TrimSpace(result), "\n")

	// 验证只有一行输出
	if len(lines) != 1 {
		t.Errorf("expected 1 line, got %d", len(lines))
	}

	// 验证包含了消息和堆栈的部分关键字（堆栈通常包含 testing 相关的路径）
	if !strings.Contains(lines[0], "main message") {
		t.Errorf("result should contain main message")
	}
	// 因为 CollectRecord 在 Warn 时会记录堆栈，所以 result 应该包含多余的信息
	if !strings.Contains(lines[0], "logf_test.go") && !strings.Contains(lines[0], "testing") {
		t.Errorf("result should contain stack information")
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
