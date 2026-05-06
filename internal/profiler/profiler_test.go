package profiler_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envlayer/envlayer/internal/profiler"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "profiler-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := tempDir(t)
	vars := map[string]string{"APP_ENV": "staging", "PORT": "8080"}

	if err := profiler.Save(dir, "staging", vars); err != nil {
		t.Fatalf("Save: %v", err)
	}

	p, err := profiler.Load(dir, "staging")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if p.Name != "staging" {
		t.Errorf("Name: got %q, want %q", p.Name, "staging")
	}
	if p.Vars["APP_ENV"] != "staging" {
		t.Errorf("APP_ENV: got %q, want %q", p.Vars["APP_ENV"], "staging")
	}
	if p.Vars["PORT"] != "8080" {
		t.Errorf("PORT: got %q, want %q", p.Vars["PORT"], "8080")
	}
}

func TestLoad_MissingProfile_ReturnsError(t *testing.T) {
	dir := tempDir(t)
	_, err := profiler.Load(dir, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing profile, got nil")
	}
}

func TestList_ReturnsSortedNames(t *testing.T) {
	dir := tempDir(t)
	for _, name := range []string{"prod", "dev", "staging"} {
		if err := profiler.Save(dir, name, map[string]string{"K": name}); err != nil {
			t.Fatalf("Save %s: %v", name, err)
		}
	}

	names, err := profiler.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	want := []string{"dev", "prod", "staging"}
	if len(names) != len(want) {
		t.Fatalf("List len: got %d, want %d", len(names), len(want))
	}
	for i, n := range want {
		if names[i] != n {
			t.Errorf("names[%d]: got %q, want %q", i, names[i], n)
		}
	}
}

func TestList_EmptyDir_ReturnsEmpty(t *testing.T) {
	dir := tempDir(t)
	names, err := profiler.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty list, got %v", names)
	}
}

func TestList_NonexistentDir_ReturnsEmpty(t *testing.T) {
	names, err := profiler.List(filepath.Join(os.TempDir(), "no-such-profiler-dir"))
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty list, got %v", names)
	}
}

func TestDelete_RemovesProfile(t *testing.T) {
	dir := tempDir(t)
	if err := profiler.Save(dir, "temp", map[string]string{"X": "1"}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := profiler.Delete(dir, "temp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	names, _ := profiler.List(dir)
	if len(names) != 0 {
		t.Errorf("expected empty after delete, got %v", names)
	}
}

func TestDelete_MissingProfile_ReturnsError(t *testing.T) {
	dir := tempDir(t)
	if err := profiler.Delete(dir, "ghost"); err == nil {
		t.Fatal("expected error deleting nonexistent profile, got nil")
	}
}
