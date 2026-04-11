package watcher

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

func TestWatch_DetectsFileChange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	if err := os.WriteFile(path, []byte("initial"), 0644); err != nil {
		t.Fatal(err)
	}

	var called atomic.Int32
	stop, err := Watch([]string{path}, func() {
		called.Add(1)
	})
	if err != nil {
		t.Fatal(err)
	}
	defer stop()

	// Give the watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Modify the file
	if err := os.WriteFile(path, []byte("modified"), 0644); err != nil {
		t.Fatal(err)
	}

	// Wait for callback (debounce is ~100ms, so wait up to 500ms)
	deadline := time.After(500 * time.Millisecond)
	for called.Load() == 0 {
		select {
		case <-deadline:
			t.Fatal("callback was not called within timeout")
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func TestWatch_EmptyPaths_NoOp(t *testing.T) {
	var called atomic.Int32
	stop, err := Watch(nil, func() {
		called.Add(1)
	})
	if err != nil {
		t.Fatal(err)
	}
	defer stop()

	time.Sleep(200 * time.Millisecond)
	if called.Load() != 0 {
		t.Fatal("callback should not have been called for empty paths")
	}
}
