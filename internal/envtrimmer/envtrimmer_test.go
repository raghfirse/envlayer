package envtrimmer_test

import (
	"testing"

	"github.com/yourusername/envlayer/internal/envtrimmer"
)

func base() map[string]string {
	return map[string]string{
		"APP_NAME":    "myapp",
		"APP_SECRET":  "supersecret",
		"DEBUG":       "",
		"INTERNAL_ID": "42",
		"LOG_LEVEL":   "info",
	}
}

func TestTrim_OmitEmpty(t *testing.T) {
	opts := envtrimmer.DefaultOptions()
	opts.OmitEmpty = true

	out := envtrimmer.Trim(base(), opts)
	if _, ok := out["DEBUG"]; ok {
		t.Error("expected DEBUG to be omitted")
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME=myapp, got %q", out["APP_NAME"])
	}
}

func TestTrim_OmitKeys(t *testing.T) {
	opts := envtrimmer.DefaultOptions()
	opts.OmitKeys = []string{"APP_SECRET", "INTERNAL_ID"}

	out := envtrimmer.Trim(base(), opts)
	if _, ok := out["APP_SECRET"]; ok {
		t.Error("expected APP_SECRET to be omitted")
	}
	if _, ok := out["INTERNAL_ID"]; ok {
		t.Error("expected INTERNAL_ID to be omitted")
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME to survive, got %q", out["APP_NAME"])
	}
}

func TestTrim_OmitPrefixes(t *testing.T) {
	opts := envtrimmer.DefaultOptions()
	opts.OmitPrefixes = []string{"APP_", "INTERNAL_"}

	out := envtrimmer.Trim(base(), opts)
	for _, k := range []string{"APP_NAME", "APP_SECRET", "INTERNAL_ID"} {
		if _, ok := out[k]; ok {
			t.Errorf("expected %s to be omitted by prefix rule", k)
		}
	}
	if out["LOG_LEVEL"] != "info" {
		t.Error("expected LOG_LEVEL to survive")
	}
}

func TestTrim_MaxValueLen(t *testing.T) {
	opts := envtrimmer.DefaultOptions()
	opts.MaxValueLen = 4

	out := envtrimmer.Trim(base(), opts)
	if out["APP_SECRET"] != "supe" {
		t.Errorf("expected value capped to 4 chars, got %q", out["APP_SECRET"])
	}
	if out["APP_NAME"] != "myap" {
		t.Errorf("expected APP_NAME capped, got %q", out["APP_NAME"])
	}
}

func TestTrim_DoesNotMutateInput(t *testing.T) {
	input := base()
	opts := envtrimmer.DefaultOptions()
	opts.OmitEmpty = true
	opts.MaxValueLen = 3

	_ = envtrimmer.Trim(input, opts)

	if input["APP_SECRET"] != "supersecret" {
		t.Error("Trim must not mutate the input map")
	}
	if _, ok := input["DEBUG"]; !ok {
		t.Error("Trim must not delete keys from the input map")
	}
}

func TestTrimKeys_KeepsOnlyListed(t *testing.T) {
	out := envtrimmer.TrimKeys(base(), []string{"APP_NAME", "LOG_LEVEL"})
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("unexpected value for APP_NAME: %q", out["APP_NAME"])
	}
	if out["LOG_LEVEL"] != "info" {
		t.Errorf("unexpected value for LOG_LEVEL: %q", out["LOG_LEVEL"])
	}
}

func TestTrimKeys_MissingKeySkipped(t *testing.T) {
	out := envtrimmer.TrimKeys(base(), []string{"NONEXISTENT"})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d entries", len(out))
	}
}
