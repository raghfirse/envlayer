package envexpander_test

import (
	"testing"

	"github.com/user/envlayer/internal/envexpander"
)

var sampleExpansions = map[string]string{
	"DB_HOST": "DATABASE_HOST",
	"DB_PORT": "DATABASE_PORT",
	"DB_PASS": "DATABASE_PASSWORD",
}

func TestExpand_ExpandsKnownKeys(t *testing.T) {
	vars := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	got, err := envexpander.Expand(vars, sampleExpansions, envexpander.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", got["DATABASE_HOST"])
	}
	if got["DATABASE_PORT"] != "5432" {
		t.Errorf("expected DATABASE_PORT=5432, got %q", got["DATABASE_PORT"])
	}
	if _, ok := got["DB_HOST"]; ok {
		t.Error("original key DB_HOST should not be present")
	}
}

func TestExpand_UnknownKeyKeptByDefault(t *testing.T) {
	vars := map[string]string{"UNKNOWN_KEY": "value"}
	got, err := envexpander.Expand(vars, sampleExpansions, envexpander.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["UNKNOWN_KEY"] != "value" {
		t.Errorf("expected UNKNOWN_KEY to be kept, got %v", got)
	}
}

func TestExpand_FailOnUnknown_ReturnsError(t *testing.T) {
	vars := map[string]string{"MYSTERY": "42"}
	opts := envexpander.Options{FailOnUnknown: true}
	_, err := envexpander.Expand(vars, sampleExpansions, opts)
	if err == nil {
		t.Fatal("expected error for unknown key, got nil")
	}
}

func TestExpand_KeepOriginal_RetainsBothKeys(t *testing.T) {
	vars := map[string]string{"DB_PASS": "secret"}
	opts := envexpander.Options{KeepOriginal: true}
	got, err := envexpander.Expand(vars, sampleExpansions, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DATABASE_PASSWORD"] != "secret" {
		t.Errorf("expected canonical key DATABASE_PASSWORD, got %v", got)
	}
	if got["DB_PASS"] != "secret" {
		t.Errorf("expected original key DB_PASS to be retained, got %v", got)
	}
}

func TestExpand_EmptyVars_ReturnsEmpty(t *testing.T) {
	got, err := envexpander.Expand(map[string]string{}, sampleExpansions, envexpander.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestInvert_ReversesMap(t *testing.T) {
	inv := envexpander.Invert(sampleExpansions)
	if inv["DATABASE_HOST"] != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %q", inv["DATABASE_HOST"])
	}
	if inv["DATABASE_PORT"] != "DB_PORT" {
		t.Errorf("expected DB_PORT, got %q", inv["DATABASE_PORT"])
	}
}
