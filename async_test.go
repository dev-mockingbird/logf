package logf_test

import (
	"sync"
	"testing"

	"github.com/dev-mockingbird/logf"
)

func TestAsyncLogger(t *testing.T) {
	logger := logf.Async(logf.New(), 10)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			logger.Logf(logf.Info, "hello world: %d", i)
			logger.Logf(logf.Info, "hello world: %d", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	logger.Wait()
}
