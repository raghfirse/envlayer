package envmapper_test

import (
	"testing"

	"github.com/your-org/envlayer/internal/envmapper"
)

func base() map[string]string {
	return map[string]string{
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"APP_NAME": "myapp",
	}
}

func TestApply_RemapsKey(t *testing.T) {
	opts := envmapper.Options{
		Rules: []envmapper.Rule{{From: "DB_HOST", To: "DATABASE_HOST"}},
	}
	out, err := envmapper.Apply(base(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", out["DATABASE_HOST"])
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be removed after remapping")
	}
}

func TestApply_KeepOriginal(t *testing.T) {
	opts := envmapper.Options{
		Rules: []envmapper.Rule{{From: "DB_PORT", To: "DATABASE_PORT", Keep: true}},
	}
	out, err := envmapper.Apply(base(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DATABASE_PORT"] != "5432" {
		t.Errorf("expected DATABASE_PORT=5432, got %q", out["DATABASE_PORT"])
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("expected original DB_PORT to be kept, got %q", out["DB_PORT"])
	}
}

func TestApply_DropUnmapped(t *testing.T) {
	opts := envmapper.Options{
		Rules:        []envmapper.Rule{{From: "DB_HOST", To: "DATABASE_HOST"}},
		DropUnmapped: true,
	}
	out, err := envmapper.Apply(base(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d: %v", len(out), out)
	}
}

func TestApply_MissingKey_IgnoredByDefault(t *testing.T) {
	opts := envmapper.Options{
		Rules: []envmapper.Rule{{From: "NONEXISTENT", To: "TARGET"}},
	}
	_, err := envmapper.Apply(base(), opts)
	if err != nil {
		t.Fatalf("expected no error for missing key, got: %v", err)
	}
}

func TestApply_MissingKey_FailOnMissing(t *testing.T) {
	opts := envmapper.Options{
		Rules:         []envmapper.Rule{{From: "NONEXISTENT", To: "TARGET"}},
		FailOnMissing: true,
	}
	_, err := envmapper.Apply(base(), opts)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestApply_EmptyRuleReturnsError(t *testing.T) {
	opts := envmapper.Options{
		Rules: []envmapper.Rule{{From: "", To: "TARGET"}},
	}
	_, err := envmapper.Apply(base(), opts)
	if err == nil {
		t.Fatal("expected error for empty From field")
	}
}

func TestKeys_ReturnsSortedDestinations(t *testing.T) {
	opts := envmapper.Options{
		Rules: []envmapper.Rule{
			{From: "A", To: "Z_KEY"},
			{From: "B", To: "A_KEY"},
			{From: "C", To: "M_KEY"},
		},
	}
	keys := envmapper.Keys(opts)
	want := []string{"A_KEY", "M_KEY", "Z_KEY"}
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("keys[%d] = %q, want %q", i, k, want[i])
		}
	}
}
