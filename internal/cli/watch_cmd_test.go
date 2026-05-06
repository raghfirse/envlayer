package cli_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envlayer/internal/cli"
)

func makeWatchDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestRunWatch_InvalidDir_ReturnsError(t *testing.T) {
	err := cli.RunWatch(cli.WatchOptions{
		Dir:      "/nonexistent/path",
		Env:      "",
		Interval: 50 * time.Millisecond,
		Quiet:    true,
	})
	if err == nil {
		t.Fatal("expected error for missing directory, got nil")
	}
}

func TestRunWatch_PrintsInitialState(t *testing.T) {
	dir := makeWatchDir(t, map[string]string{
		".env": "GREETING=hello\nNAME=world\n",
	})

	// RunWatch blocks on the event loop; we test only the resolver/loader
	// path by confirming no error is returned when files exist.
	// A full integration test would require stdout capture + goroutine.
	// Here we verify the option struct wires correctly through resolver.
	opts := cli.WatchOptions{
		Dir:      dir,
		Env:      "",
		Interval: 50 * time.Millisecond,
		Quiet:    true,
	}

	done := make(chan error, 1)
	go func() {
		// RunWatch will block; we just let it start and then check it
		// didn't return an error immediately.
		done <- cli.RunWatch(opts)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		// still running — expected for a watcher
	}
}

func TestRunWatch_MissingEnvFile_ReturnsError(t *testing.T) {
	dir := t.TempDir() // empty, no .env files

	err := cli.RunWatch(cli.WatchOptions{
		Dir:      dir,
		Env:      "",
		Interval: 50 * time.Millisecond,
		Quiet:    true,
	})
	if err == nil {
		t.Fatal("expected error when no env files found")
	}
}
