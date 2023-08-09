package logf

import (
	"runtime"
	"sync"
)

type asyncItem struct {
	level  Level
	format string
	args   []any
}

type async struct {
	underlying Logger
	q          chan asyncItem
	ql         int
	wg         sync.WaitGroup
	stopCh     chan struct{}
}

type AsyncLogger interface {
	Logger
	Wait()
}

func Async(underlying Logger, ql int) AsyncLogger {
	a := &async{
		underlying: underlying,
		q:          make(chan asyncItem, ql),
		ql:         ql,
		stopCh:     make(chan struct{}),
	}
	go func() {
		a.start()
	}()
	runtime.SetFinalizer(a, (*async).stop)
	return a
}

func (a *async) Prefix(prefix string) Logger {
	return Async(a.underlying.Prefix(prefix), a.ql)
}

func (a *async) Logf(l Level, format string, args ...any) {
	a.q <- asyncItem{level: l, format: format, args: args}
	a.wg.Add(1)
}

func (a *async) start() error {
	for {
		select {
		case item := <-a.q:
			a.underlying.Logf(item.level, item.format, item.args...)
			a.wg.Done()
		case <-a.stopCh:
			close(a.q)
			return nil
		}
	}
}

func (a *async) Wait() {
	a.wg.Wait()
}

func (a *async) stop() error {
	runtime.SetFinalizer(a, nil)
	a.stopCh <- struct{}{}
	return nil
}
