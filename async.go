package logf

import (
	"runtime"
	"sync"
)

type asyncLogger struct {
	config *Config
	lock   sync.Mutex
	q      chan Record
	wg     sync.WaitGroup
	stopCh chan struct{}
}

type AsyncLogger interface {
	Logger
	Wait()
}

func Async(opts ...Option) AsyncLogger {
	a := &asyncLogger{
		config: getConfig(opts...),
	}
	a.q = make(chan Record, a.config.BufferSize)
	a.stopCh = make(chan struct{})
	go func() {
		a.start()
	}()
	runtime.SetFinalizer(a, (*asyncLogger).stop)
	return a
}

func (a *asyncLogger) Prefix(prefix string) Logger {
	ret := &asyncLogger{
		config: a.config.WithPrefix(prefix),
		q:      make(chan Record, a.config.BufferSize),
		stopCh: make(chan struct{}),
	}
	go func() {
		ret.start()
	}()
	runtime.SetFinalizer(ret, (*asyncLogger).stop)
	return ret
}

func (a *asyncLogger) Logf(l Level, format string, args ...any) {
	if l < a.config.LogLevel {
		return
	}
	record := CollectRecord(l, a.config.CallerDepth, format, args...)
	a.wg.Add(1)
	a.q <- record
}

func (a *asyncLogger) start() error {
	for {
		select {
		case item := <-a.q:
			a.lock.Lock()
			a.config.Printer.Print(a.config.Prefix, item)
			a.lock.Unlock()
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
