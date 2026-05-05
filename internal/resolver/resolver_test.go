package resolver_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envlayer/internal/resolver"
)

func makeFiles(t *testing.T, dir string, names []string) {
	t.Helper()
	for _, name := range names {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte("KEY=val\n"), 0644); err != nil {
			t.Fatalf("failed to create file %s: %v", path, err)
		}
	}
}

func TestResolve_BaseOnly(t *testing.T) {
	dir := t.TempDir()
	makeFiles(t, dir, []string{".env"})

	paths, err := resolver.Resolve(resolver.Config{BaseDir: dir, Environment: "development"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 1 || filepath.Base(paths[0]) != ".env" {
		t.Errorf("expected only .env, got %v", paths)
	}
}

func TestResolve_WithEnvironmentOverride(t *testing.T) {
	dir := t.TempDir()
	makeFiles(t, dir, []string{".env", ".env.staging"})

	paths, err := resolver.Resolve(resolver.Config{BaseDir: dir, Environment: "staging"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %d: %v", len(paths), paths)
	}
	if filepath.Base(paths[0]) != ".env" {
		t.Errorf("expected .env first, got %s", paths[0])
	}
	if filepath.Base(paths[1]) != ".env.staging" {
		t.Errorf("expected .env.staging second, got %s", paths[1])
	}
}

func TestResolve_LocalOverrideIncluded(t *testing.T) {
	dir := t.TempDir()
	makeFiles(t, dir, []string{".env", ".env.production", ".env.production.local"})

	paths, err := resolver.Resolve(resolver.Config{BaseDir: dir, Environment: "production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 3 {
		t.Fatalf("expected 3 paths, got %d: %v", len(paths), paths)
	}
}

func TestResolve_NoFilesReturnsError(t *testing.T) {
	dir := t.TempDir()

	_, err := resolver.Resolve(resolver.Config{BaseDir: dir, Environment: "test"})
	if err == nil {
		t.Fatal("expected error when no .env files exist")
	}
}

func TestResolve_DefaultBaseDir(t *testing.T) {
	// Should not panic or crash with empty BaseDir.
	_, err := resolver.Resolve(resolver.Config{Environment: "dev"})
	// We don't assert success here since CWD may not have .env files.
	_ = err
}
