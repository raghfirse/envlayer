package envaliaser_test

import (
	"testing"

	"github.com/yourorg/envlayer/internal/envaliaser"
)

func TestApply_RenamesKey(t *testing.T) {
	vars := map[string]string{"DB_HOST": "localhost"}
	aliases := envaliaser.AliasMap{"DB_HOST": {"DATABASE_HOST"}}
	out, err := envaliaser.Apply(vars, aliases, envaliaser.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", out["DATABASE_HOST"])
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("original key DB_HOST should have been removed")
	}
}

func TestApply_KeepOriginal(t *testing.T) {
	vars := map[string]string{"API_KEY": "secret"}
	aliases := envaliaser.AliasMap{"API_KEY": {"SERVICE_API_KEY"}}
	opts := envaliaser.Options{KeepOriginal: true}
	out, err := envaliaser.Apply(vars, aliases, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_KEY"] != "secret" {
		t.Error("original key should be retained")
	}
	if out["SERVICE_API_KEY"] != "secret" {
		t.Error("alias should also be set")
	}
}

func TestApply_FanOut_MultipleAliases(t *testing.T) {
	vars := map[string]string{"PORT": "8080"}
	aliases := envaliaser.AliasMap{"PORT": {"HTTP_PORT", "APP_PORT"}}
	out, err := envaliaser.Apply(vars, aliases, envaliaser.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, k := range []string{"HTTP_PORT", "APP_PORT"} {
		if out[k] != "8080" {
			t.Errorf("expected %s=8080, got %q", k, out[k])
		}
	}
}

func TestApply_MissingKey_IgnoredByDefault(t *testing.T) {
	vars := map[string]string{"EXISTING": "yes"}
	aliases := envaliaser.AliasMap{"MISSING": {"ALIAS"}}
	out, err := envaliaser.Apply(vars, aliases, envaliaser.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["ALIAS"]; ok {
		t.Error("alias for missing key should not appear in output")
	}
}

func TestApply_MissingKey_FailOnMissing(t *testing.T) {
	vars := map[string]string{}
	aliases := envaliaser.AliasMap{"GONE": {"NEW"}}
	opts := envaliaser.Options{FailOnMissing: true}
	_, err := envaliaser.Apply(vars, aliases, opts)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	vars := map[string]string{"X": "1"}
	original := map[string]string{"X": "1"}
	aliases := envaliaser.AliasMap{"X": {"Y"}}
	_, _ = envaliaser.Apply(vars, aliases, envaliaser.DefaultOptions())
	if vars["X"] != original["X"] {
		t.Error("input map was mutated")
	}
}

func TestInvert_ReversesMapping(t *testing.T) {
	aliases := envaliaser.AliasMap{
		"DB_HOST": {"DATABASE_HOST", "PG_HOST"},
	}
	inv := envaliaser.Invert(aliases)
	if inv["DATABASE_HOST"] != "DB_HOST" {
		t.Errorf("expected DATABASE_HOST -> DB_HOST, got %q", inv["DATABASE_HOST"])
	}
	if inv["PG_HOST"] != "DB_HOST" {
		t.Errorf("expected PG_HOST -> DB_HOST, got %q", inv["PG_HOST"])
	}
}
