package logf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type DailyRotateWriter struct {
	dir     string
	prefix  string
	maxAge  int // days
	mu      sync.Mutex
	current *os.File
	curDate string
}

func NewDailyRotateWriter(dir, prefix string, maxAgeDays int) (*DailyRotateWriter, error) {
	if maxAgeDays <= 0 {
		maxAgeDays = 30
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create log dir: %w", err)
	}
	w := &DailyRotateWriter{dir: dir, prefix: prefix, maxAge: maxAgeDays}
	if err := w.rotate(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *DailyRotateWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	today := time.Now().Format("2006-01-02")
	if today != w.curDate {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}
	return w.current.Write(p)
}

func (w *DailyRotateWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.current != nil {
		return w.current.Close()
	}
	return nil
}

func (w *DailyRotateWriter) rotate() error {
	today := time.Now().Format("2006-01-02")
	currentPath := filepath.Join(w.dir, fmt.Sprintf("%s.log", w.prefix))
	if w.current != nil {
		w.current.Close()
		// rename the old file to a dated archive
		if w.curDate != "" && w.curDate != today {
			archived := filepath.Join(w.dir, fmt.Sprintf("%s-%s.log", w.prefix, w.curDate))
			os.Rename(currentPath, archived)
		}
	}
	f, err := os.OpenFile(currentPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("open log file %s: %w", currentPath, err)
	}
	w.current = f
	w.curDate = today
	go w.cleanup()
	return nil
}

func (w *DailyRotateWriter) cleanup() {
	cutoff := time.Now().AddDate(0, 0, -w.maxAge)
	entries, err := os.ReadDir(w.dir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasPrefix(name, w.prefix) || !strings.HasSuffix(name, ".log") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			os.Remove(filepath.Join(w.dir, name))
		}
	}
}
