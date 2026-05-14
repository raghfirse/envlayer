package envhistory_test

import (
	"os"
	"testing"
	"time"

	"github.com/nicholasgasior/envlayer/internal/envhistory"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "envhistory-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestRecord_CreatesEntry(t *testing.T) {
	dir := tempDir(t)
	vars := map[string]string{"APP_ENV": "production", "PORT": "8080"}

	e, err := envhistory.Record(dir, "initial", vars)
	if err != nil {
		t.Fatalf("Record: %v", err)
	}
	if e.Label != "initial" {
		t.Errorf("label = %q, want %q", e.Label, "initial")
	}
	if e.Vars["APP_ENV"] != "production" {
		t.Errorf("APP_ENV = %q, want production", e.Vars["APP_ENV"])
	}
	if e.ID == "" {
		t.Error("expected non-empty ID")
	}
}

func TestRecord_DoesNotMutateInput(t *testing.T) {
	dir := tempDir(t)
	vars := map[string]string{"KEY": "value"}

	_, _ = envhistory.Record(dir, "test", vars)
	vars["KEY"] = "mutated"

	entries, _ := envhistory.List(dir)
	if entries[0].Vars["KEY"] != "value" {
		t.Errorf("expected stored value to be unchanged, got %q", entries[0].Vars["KEY"])
	}
}

func TestList_ReturnsSortedByTime(t *testing.T) {
	dir := tempDir(t)

	_, _ = envhistory.Record(dir, "first", map[string]string{"A": "1"})
	time.Sleep(2 * time.Millisecond)
	_, _ = envhistory.Record(dir, "second", map[string]string{"B": "2"})

	entries, err := envhistory.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Label != "first" {
		t.Errorf("expected first entry label %q, got %q", "first", entries[0].Label)
	}
	if entries[1].Label != "second" {
		t.Errorf("expected second entry label %q, got %q", "second", entries[1].Label)
	}
}

func TestList_EmptyDir_ReturnsNil(t *testing.T) {
	dir := tempDir(t)
	entries, err := envhistory.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestGet_ReturnsEntry(t *testing.T) {
	dir := tempDir(t)
	e, _ := envhistory.Record(dir, "snap", map[string]string{"X": "42"})

	got, err := envhistory.Get(dir, e.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Label != "snap" {
		t.Errorf("label = %q, want snap", got.Label)
	}
}

func TestGet_MissingEntry_ReturnsError(t *testing.T) {
	dir := tempDir(t)
	_, err := envhistory.Get(dir, "nonexistent")
	if err == nil {
		t.Error("expected error for missing entry")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	dir := tempDir(t)
	e, _ := envhistory.Record(dir, "to-delete", map[string]string{"K": "v"})

	if err := envhistory.Delete(dir, e.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	entries, _ := envhistory.List(dir)
	if len(entries) != 0 {
		t.Errorf("expected 0 entries after delete, got %d", len(entries))
	}
}

func TestDelete_MissingEntry_ReturnsError(t *testing.T) {
	dir := tempDir(t)
	if err := envhistory.Delete(dir, "ghost"); err == nil {
		t.Error("expected error for missing entry")
	}
}
