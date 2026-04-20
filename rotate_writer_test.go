package logf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestDailyRotateWriter_Write(t *testing.T) {
	dir := t.TempDir()
	w, err := NewDailyRotateWriter(dir, "app", 30)
	if err != nil {
		t.Fatalf("NewDailyRotateWriter: %v", err)
	}
	defer w.Close()

	msg := "hello rotate\n"
	n, err := w.Write([]byte(msg))
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	if n != len(msg) {
		t.Errorf("wrote %d bytes, want %d", n, len(msg))
	}

	// current file should be app.log (no date)
	currentPath := filepath.Join(dir, "app.log")
	data, err := os.ReadFile(currentPath)
	if err != nil {
		t.Fatalf("read current file: %v", err)
	}
	if !strings.Contains(string(data), "hello rotate") {
		t.Errorf("current file does not contain written message")
	}
}

func TestDailyRotateWriter_CurrentFileHasNoDate(t *testing.T) {
	dir := t.TempDir()
	w, err := NewDailyRotateWriter(dir, "app", 30)
	if err != nil {
		t.Fatalf("NewDailyRotateWriter: %v", err)
	}
	defer w.Close()

	w.Write([]byte("data\n"))

	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		name := e.Name()
		// the only file should be app.log, not app-2006-01-02.log
		today := time.Now().Format("2006-01-02")
		if strings.Contains(name, today) {
			t.Errorf("current file should not contain date, got: %s", name)
		}
	}
}

func TestDailyRotateWriter_RotateArchivesWithDate(t *testing.T) {
	dir := t.TempDir()
	w, err := NewDailyRotateWriter(dir, "app", 30)
	if err != nil {
		t.Fatalf("NewDailyRotateWriter: %v", err)
	}
	defer w.Close()

	w.Write([]byte("before rotate\n"))

	// simulate a date change by manually setting curDate to yesterday
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	w.mu.Lock()
	w.curDate = yesterday
	w.mu.Unlock()

	// trigger rotate via Write
	w.Write([]byte("after rotate\n"))

	// archived file should exist: app-<yesterday>.log
	archived := filepath.Join(dir, fmt.Sprintf("app-%s.log", yesterday))
	if _, err := os.Stat(archived); os.IsNotExist(err) {
		t.Errorf("expected archived file %s to exist", archived)
	}

	// current file app.log should exist and contain new message
	data, err := os.ReadFile(filepath.Join(dir, "app.log"))
	if err != nil {
		t.Fatalf("read current file: %v", err)
	}
	if !strings.Contains(string(data), "after rotate") {
		t.Errorf("current file should contain 'after rotate'")
	}
}

func TestDailyRotateWriter_CleanupOldFiles(t *testing.T) {
	dir := t.TempDir()

	// create an old dated file (older than maxAge)
	oldDate := time.Now().AddDate(0, 0, -5).Format("2006-01-02")
	oldFile := filepath.Join(dir, fmt.Sprintf("app-%s.log", oldDate))
	os.WriteFile(oldFile, []byte("old\n"), 0644)
	// backdate its mtime
	old := time.Now().AddDate(0, 0, -5)
	os.Chtimes(oldFile, old, old)

	w, err := NewDailyRotateWriter(dir, "app", 3)
	if err != nil {
		t.Fatalf("NewDailyRotateWriter: %v", err)
	}
	defer w.Close()

	// give cleanup goroutine time to run
	time.Sleep(100 * time.Millisecond)

	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Errorf("expected old file %s to be cleaned up", oldFile)
	}
}

func TestDailyRotateWriter_ConcurrentWrite(t *testing.T) {
	dir := t.TempDir()
	w, err := NewDailyRotateWriter(dir, "app", 30)
	if err != nil {
		t.Fatalf("NewDailyRotateWriter: %v", err)
	}
	defer w.Close()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			w.Write([]byte(fmt.Sprintf("goroutine %d\n", i)))
		}(i)
	}
	wg.Wait()
}
