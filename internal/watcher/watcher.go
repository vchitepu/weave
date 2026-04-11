package watcher

import (
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watch watches the given file paths for changes and calls onChange when any
// file is written or created. Changes are debounced by 100ms.
// Returns a stop function to clean up the watcher.
// If paths is empty, returns a no-op stop function and nil error.
func Watch(paths []string, onChange func()) (stop func(), err error) {
	if len(paths) == 0 {
		return func() {}, nil
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	for _, p := range paths {
		if err := w.Add(p); err != nil {
			w.Close()
			return nil, err
		}
	}

	var once sync.Once
	done := make(chan struct{})

	go func() {
		var timer *time.Timer
		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					return
				}
				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
					if timer != nil {
						timer.Stop()
					}
					timer = time.AfterFunc(100*time.Millisecond, onChange)
				}
			case _, ok := <-w.Errors:
				if !ok {
					return
				}
				// Ignore errors silently
			case <-done:
				return
			}
		}
	}()

	stopFn := func() {
		once.Do(func() {
			close(done)
			w.Close()
		})
	}

	return stopFn, nil
}
