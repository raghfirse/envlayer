package watcher_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envlayer/internal/watcher"
)

func writeTmp(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestWatch_DetectsModification(t *testing.T) {
	dir := t.TempDir()
	p := writeTmp(t, dir, ".env", "KEY=original\n")

	done := make(chan struct{})
	defer close(done)

	events := watcher.Watch([]string{p}, 20*time.Millisecond, done)

	time.Sleep(40 * time.Millisecond)
	if err := os.WriteFile(p, []byte("KEY=changed\n"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case ev := <-events:
		if ev.Kind != watcher.ChangeModified {
			t.Fatalf("expected modified, got %s", ev.Kind)
		}
		if ev.Path != p {
			t.Fatalf("unexpected path %s", ev.Path)
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for modification event")
	}
}

func TestWatch_DetectsRemoval(t *testing.T) {
	dir := t.TempDir()
	p := writeTmp(t, dir, ".env", "KEY=value\n")

	done := make(chan struct{})
	defer close(done)

	events := watcher.Watch([]string{p}, 20*time.Millisecond, done)

	time.Sleep(40 * time.Millisecond)
	os.Remove(p)

	select {
	case ev := <-events:
		if ev.Kind != watcher.ChangeRemoved {
			t.Fatalf("expected removed, got %s", ev.Kind)
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for removal event")
	}
}

func TestWatch_DetectsCreation(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env.local")

	done := make(chan struct{})
	defer close(done)

	events := watcher.Watch([]string{p}, 20*time.Millisecond, done)

	time.Sleep(40 * time.Millisecond)
	if err := os.WriteFile(p, []byte("NEW=1\n"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case ev := <-events:
		if ev.Kind != watcher.ChangeCreated {
			t.Fatalf("expected created, got %s", ev.Kind)
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for creation event")
	}
}

func TestWatch_NoEventWhenUnchanged(t *testing.T) {
	dir := t.TempDir()
	p := writeTmp(t, dir, ".env", "STABLE=yes\n")

	done := make(chan struct{})
	defer close(done)

	events := watcher.Watch([]string{p}, 20*time.Millisecond, done)

	select {
	case ev := <-events:
		t.Fatalf("unexpected event: %+v", ev)
	case <-time.After(150 * time.Millisecond):
		// pass — no spurious events
	}
}
