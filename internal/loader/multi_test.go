package loader_test

import (
	"path/filepath"
	"testing"

	"github.com/user/envlayer/internal/loader"
)

func TestLoadFiles_MergesInOrder(t *testing.T) {
	dir := t.TempDir()

	base := filepath.Join(dir, ".env")
	override := filepath.Join(dir, ".env.production")

	writeTempEnv(t, base, "APP=base\nDEBUG=false\n")
	writeTempEnv(t, override, "DEBUG=true\nNEW_KEY=hello\n")

	result, err := loader.LoadFiles([]string{base, override})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["APP"] != "base" {
		t.Errorf("expected APP=base, got %q", result["APP"])
	}
	if result["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true (overridden), got %q", result["DEBUG"])
	}
	if result["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q", result["NEW_KEY"])
	}
}

func TestLoadFiles_SingleFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	writeTempEnv(t, path, "ONLY=one\n")

	result, err := loader.LoadFiles([]string{path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["ONLY"] != "one" {
		t.Errorf("expected ONLY=one, got %q", result["ONLY"])
	}
}

func TestLoadFiles_EmptyList(t *testing.T) {
	result, err := loader.LoadFiles([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestLoadFiles_MissingFileReturnsError(t *testing.T) {
	_, err := loader.LoadFiles([]string{"/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFiles_LaterFileOverridesAllPreviousValues(t *testing.T) {
	dir := t.TempDir()

	first := filepath.Join(dir, ".env")
	second := filepath.Join(dir, ".env.local")
	third := filepath.Join(dir, ".env.production")

	writeTempEnv(t, first, "KEY=first\n")
	writeTempEnv(t, second, "KEY=second\n")
	writeTempEnv(t, third, "KEY=third\n")

	result, err := loader.LoadFiles([]string{first, second, third})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "third" {
		t.Errorf("expected KEY=third (last file wins), got %q", result["KEY"])
	}
}
