package logf

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

type AsyncLogRecord struct {
	CreatedAt time.Time
	PathFile  string
	Line      int
	Stack     [][]byte
	Level     Level
	Format    string
	Args      []any
}

type asyncLogPrinter interface {
	Print(prefix string, record AsyncLogRecord)
}

type asyncPrinter struct {
	w io.Writer
}

func NewAsyncPrinter(w io.Writer) asyncPrinter {
	return asyncPrinter{w: w}
}

func (a asyncPrinter) Print(prefix string, record AsyncLogRecord) {
	fmt.Fprintf(
		a.w,
		"%s %s:%d:\t[%s]\t%s%s\n",
		record.CreatedAt.Format("2006/01/02 15:04:05"),
		path.Base(record.PathFile),
		record.Line,
		LevelString(record.Level),
		prefix,
		fmt.Sprintf(record.Format, record.Args...))
	if record.Level >= Warn {
		fmt.Fprintf(a.w, "%s\n", string(bytes.Join(record.Stack, []byte{'\n'})))
	}
}

type asyncLogger struct {
	printer  asyncLogPrinter
	prefix   string
	q        chan AsyncLogRecord
	ql       int
	wg       sync.WaitGroup
	logLevel Level
	stopCh   chan struct{}
}

type AsyncLogger interface {
	Logger
	Wait()
}

type AsyncLoggerOption func(l *asyncLogger)

func AsyncPrinter(p asyncLogPrinter) AsyncLoggerOption {
	return func(l *asyncLogger) {
		l.printer = p
	}
}

func AsyncLevel(l Level) AsyncLoggerOption {
	return func(logger *asyncLogger) {
		logger.logLevel = l
	}
}

func Async(opts ...AsyncLoggerOption) AsyncLogger {
	a := &asyncLogger{
		logLevel: Info,
		printer:  NewAsyncPrinter(os.Stdout),
		q:        make(chan AsyncLogRecord, 100),
		ql:       100,
		stopCh:   make(chan struct{}),
	}
	go func() {
		a.start()
	}()
	runtime.SetFinalizer(a, (*asyncLogger).stop)
	return a
}

func (a *asyncLogger) Prefix(prefix string) Logger {
	return &asyncLogger{
		printer: a.printer,
		q:       make(chan AsyncLogRecord, a.ql),
		ql:      a.ql,
		stopCh:  make(chan struct{}),
	}
}

func (a *asyncLogger) Logf(l Level, format string, args ...any) {
	record := AsyncLogRecord{CreatedAt: time.Now(), Level: l, Format: format, Args: args}
	if _, file, line, ok := runtime.Caller(1); ok {
		record.PathFile = file
		record.Line = line
	}
	if l >= Warn {
		stack := debug.Stack()
		record.Stack = bytes.Split(stack, []byte{'\n'})[5:]
	}
	a.q <- record
	a.wg.Add(1)
}

func (a *asyncLogger) start() error {
	for {
		select {
		case item := <-a.q:
			a.printer.Print(a.prefix, item)
			a.wg.Done()
		case <-a.stopCh:
			close(a.q)
			return nil
		}
	}
}

func (a *asyncLogger) Wait() {
	a.wg.Wait()
}

func (a *asyncLogger) stop() error {
	runtime.SetFinalizer(a, nil)
	a.stopCh <- struct{}{}
	return nil
}
