package envrewriter_test

import (
	"testing"

	"github.com/user/envlayer/internal/envrewriter"
)

func TestRewrite_ReplaceInValues(t *testing.T) {
	vars := map[string]string{"DB_HOST": "localhost", "APP_HOST": "localhost"}
	opts := envrewriter.DefaultOptions()
	opts.Rules = []envrewriter.Rule{
		{Target: "value", Find: "localhost", Replace: "db.internal"},
	}
	out, err := envrewriter.Rewrite(vars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "db.internal" {
		t.Errorf("expected db.internal, got %s", out["DB_HOST"])
	}
	if out["APP_HOST"] != "db.internal" {
		t.Errorf("expected db.internal, got %s", out["APP_HOST"])
	}
}

func TestRewrite_ReplaceInKeys(t *testing.T) {
	vars := map[string]string{"OLD_TOKEN": "abc", "OLD_SECRET": "xyz"}
	opts := envrewriter.DefaultOptions()
	opts.Rules = []envrewriter.Rule{
		{Target: "key", Find: "OLD_", Replace: "NEW_"},
	}
	out, err := envrewriter.Rewrite(vars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["NEW_TOKEN"]; !ok {
		t.Error("expected NEW_TOKEN key")
	}
	if _, ok := out["NEW_SECRET"]; !ok {
		t.Error("expected NEW_SECRET key")
	}
	if _, ok := out["OLD_TOKEN"]; ok {
		t.Error("OLD_TOKEN should be gone")
	}
}

func TestRewrite_ReplaceInBoth(t *testing.T) {
	vars := map[string]string{"STAGE_URL": "http://stage.example.com"}
	opts := envrewriter.DefaultOptions()
	opts.Rules = []envrewriter.Rule{
		{Target: "both", Find: "stage", Replace: "prod"},
	}
	out, err := envrewriter.Rewrite(vars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := out["PROD_URL"]; !ok || v != "http://prod.example.com" {
		t.Errorf("unexpected result: key=%v val=%v", ok, v)
	}
}

func TestRewrite_CaseInsensitive(t *testing.T) {
	vars := map[string]string{"API_KEY": "MySecretValue"}
	opts := envrewriter.Options{
		CaseSensitive: false,
		Rules: []envrewriter.Rule{
			{Target: "value", Find: "secret", Replace: "***"},
		},
	}
	out, err := envrewriter.Rewrite(vars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_KEY"] != "My***Value" {
		t.Errorf("expected My***Value, got %s", out["API_KEY"])
	}
}

func TestRewrite_EmptyFindSkipsRule(t *testing.T) {
	vars := map[string]string{"FOO": "bar"}
	opts := envrewriter.DefaultOptions()
	opts.Rules = []envrewriter.Rule{
		{Target: "value", Find: "", Replace: "replaced"},
	}
	out, err := envrewriter.Rewrite(vars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" {
		t.Errorf("expected bar, got %s", out["FOO"])
	}
}

func TestRewrite_DoesNotMutateInput(t *testing.T) {
	vars := map[string]string{"X": "original"}
	opts := envrewriter.DefaultOptions()
	opts.Rules = []envrewriter.Rule{
		{Target: "value", Find: "original", Replace: "changed"},
	}
	_, _ = envrewriter.Rewrite(vars, opts)
	if vars["X"] != "original" {
		t.Error("input map was mutated")
	}
}
