package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envlayer/internal/snapshot"
)

func TestTake_CopiesVars(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s := snapshot.Take(env, "production")

	if s.Environment != "production" {
		t.Errorf("expected environment 'production', got %q", s.Environment)
	}
	if s.Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", s.Vars["FOO"])
	}
	// Mutating original should not affect snapshot
	env["FOO"] = "mutated"
	if s.Vars["FOO"] != "bar" {
		t.Error("snapshot vars were mutated by original map change")
	}
}

func TestTake_SetsCreatedAt(t *testing.T) {
	before := time.Now().UTC()
	s := snapshot.Take(map[string]string{}, "")
	after := time.Now().UTC()

	if s.CreatedAt.Before(before) || s.CreatedAt.After(after) {
		t.Errorf("CreatedAt %v not within expected range", s.CreatedAt)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	original := snapshot.Take(map[string]string{"KEY": "value", "PORT": "8080"}, "staging")
	if err := snapshot.Save(original, path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.Environment != "staging" {
		t.Errorf("expected environment 'staging', got %q", loaded.Environment)
	}
	if loaded.Vars["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", loaded.Vars["KEY"])
	}
	if loaded.Vars["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", loaded.Vars["PORT"])
	}
}

func TestLoad_MissingFile_ReturnsError(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not-json{"), 0644)

	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestSortedKeys(t *testing.T) {
	s := snapshot.Take(map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}, "")
	keys := s.SortedKeys()
	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: expected %q, got %q", i, expected[i], k)
		}
	}
}
