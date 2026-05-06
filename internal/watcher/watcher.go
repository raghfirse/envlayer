package watcher

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"
)

// FileState holds the last known state of a watched file.
type FileState struct {
	Path    string
	ModTime time.Time
	Checksum string
}

// ChangeEvent describes a detected change in a watched file.
type ChangeEvent struct {
	Path    string
	Kind    ChangeKind
}

// ChangeKind categorises the type of file change.
type ChangeKind string

const (
	ChangeModified ChangeKind = "modified"
	ChangeCreated  ChangeKind = "created"
	ChangeRemoved  ChangeKind = "removed"
)

// Watch polls the given file paths at the specified interval and sends
// ChangeEvents on the returned channel. Cancel by closing done.
func Watch(paths []string, interval time.Duration, done <-chan struct{}) <-chan ChangeEvent {
	events := make(chan ChangeEvent, len(paths))
	states := make(map[string]FileState)

	for _, p := range paths {
		if st, err := stat(p); err == nil {
			states[p] = st
		}
	}

	go func() {
		defer close(events)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				for _, p := range paths {
					prev, known := states[p]
					curr, err := stat(p)
					if err != nil {
						if known {
							delete(states, p)
							events <- ChangeEvent{Path: p, Kind: ChangeRemoved}
						}
						continue
					}
					if !known {
						states[p] = curr
						events <- ChangeEvent{Path: p, Kind: ChangeCreated}
						continue
					}
					if curr.Checksum != prev.Checksum {
						states[p] = curr
						events <- ChangeEvent{Path: p, Kind: ChangeModified}
					}
				}
			}
		}
	}()

	return events
}

func stat(path string) (FileState, error) {
	f, err := os.Open(path)
	if err != nil {
		return FileState{}, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return FileState{}, err
	}

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return FileState{}, err
	}

	return FileState{
		Path:     path,
		ModTime:  info.ModTime(),
		Checksum: fmt.Sprintf("%x", h.Sum(nil)),
	}, nil
}
