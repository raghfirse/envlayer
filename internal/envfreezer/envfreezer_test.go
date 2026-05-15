package envfreezer_test

import (
	"testing"

	"github.com/nicholasgasior/envlayer/internal/envfreezer"
)

func TestFreeze_FrozenKeyNotOverridden(t *testing.T) {
	base := map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
	updates := map[string]string{"DB_HOST": "prod-db", "PORT": "9999", "NEW_KEY": "hello"}

	result, err := envfreezer.Freeze(base, updates, envfreezer.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST to remain %q, got %q", "localhost", result["DB_HOST"])
	}
	if result["PORT"] != "5432" {
		t.Errorf("expected PORT to remain %q, got %q", "5432", result["PORT"])
	}
	if result["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY %q, got %q", "hello", result["NEW_KEY"])
	}
}

func TestFreeze_AllowOverrideExemptsKey(t *testing.T) {
	base := map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
	updates := map[string]string{"DB_HOST": "prod-db", "PORT": "9999"}
	opts := envfreezer.Options{AllowOverride: []string{"DB_HOST"}}

	result, err := envfreezer.Freeze(base, updates, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["DB_HOST"] != "prod-db" {
		t.Errorf("expected DB_HOST %q, got %q", "prod-db", result["DB_HOST"])
	}
	if result["PORT"] != "5432" {
		t.Errorf("expected PORT to remain frozen as %q, got %q", "5432", result["PORT"])
	}
}

func TestFreeze_StrictMode_ReturnsErrorOnConflict(t *testing.T) {
	base := map[string]string{"SECRET": "abc"}
	updates := map[string]string{"SECRET": "xyz"}
	opts := envfreezer.Options{StrictMode: true}

	_, err := envfreezer.Freeze(base, updates, opts)
	if err == nil {
		t.Fatal("expected error in strict mode, got nil")
	}
}

func TestFreeze_StrictMode_SameValueNoError(t *testing.T) {
	base := map[string]string{"SECRET": "abc"}
	updates := map[string]string{"SECRET": "abc"}
	opts := envfreezer.Options{StrictMode: true}

	_, err := envfreezer.Freeze(base, updates, opts)
	if err != nil {
		t.Fatalf("unexpected error for unchanged frozen key: %v", err)
	}
}

func TestFreeze_NewKeysAlwaysAdded(t *testing.T) {
	base := map[string]string{"A": "1"}
	updates := map[string]string{"B": "2", "C": "3"}

	result, err := envfreezer.Freeze(base, updates, envfreezer.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["B"] != "2" || result["C"] != "3" {
		t.Errorf("new keys not added correctly: %v", result)
	}
}

func TestFrozenKeys_ExcludesAllowOverride(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2", "C": "3"}
	opts := envfreezer.Options{AllowOverride: []string{"B"}}

	keys := envfreezer.FrozenKeys(base, opts)
	for _, k := range keys {
		if k == "B" {
			t.Error("expected B to be excluded from frozen keys")
		}
	}
	if len(keys) != 2 {
		t.Errorf("expected 2 frozen keys, got %d", len(keys))
	}
}

func TestFrozenKeys_SortedOutput(t *testing.T) {
	base := map[string]string{"Z": "1", "A": "2", "M": "3"}
	keys := envfreezer.FrozenKeys(base, envfreezer.DefaultOptions())
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("keys not sorted: %v", keys)
		}
	}
}
