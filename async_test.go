package logf_test

import (
	"bytes"
	"fmt"
	"sync"
	"testing"

	"github.com/dev-mockingbird/logf"
)

func TestAsyncLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := logf.Async(logf.Writer(&buf))
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			logger.Logf(logf.Info, "hello world: %d", i)
			logger.Logf(logf.Warn, "hello world: %d", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Print(buf.String())
	logger.Wait()
}
