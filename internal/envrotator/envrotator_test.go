package envrotator_test

import (
	"testing"

	"github.com/your-org/envlayer/internal/envrotator"
)

func base() map[string]string {
	return map[string]string{
		"DB_PASSWORD": "old-secret",
		"API_KEY":     "key-abc",
		"DEBUG":       "true",
	}
}

func TestRotate_SingleKey(t *testing.T) {
	result, err := envrotator.Rotate(base(), []envrotator.Rotation{
		{Key: "DB_PASSWORD", NewValue: "new-secret"},
	}, envrotator.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Vars["DB_PASSWORD"] != "new-secret" {
		t.Errorf("expected new-secret, got %q", result.Vars["DB_PASSWORD"])
	}
	if len(result.Log) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(result.Log))
	}
	if result.Log[0].OldValue != "old-secret" {
		t.Errorf("expected old value old-secret, got %q", result.Log[0].OldValue)
	}
}

func TestRotate_DoesNotMutateInput(t *testing.T) {
	v := base()
	_, err := envrotator.Rotate(v, []envrotator.Rotation{
		{Key: "DB_PASSWORD", NewValue: "new-secret"},
	}, envrotator.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v["DB_PASSWORD"] != "old-secret" {
		t.Error("input map was mutated")
	}
}

func TestRotate_SkipUnchanged_NoLogEntry(t *testing.T) {
	opts := envrotator.DefaultOptions() // SkipUnchanged: true
	result, err := envrotator.Rotate(base(), []envrotator.Rotation{
		{Key: "DEBUG", NewValue: "true"}, // same value
	}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Log) != 0 {
		t.Errorf("expected no log entries for unchanged value, got %d", len(result.Log))
	}
}

func TestRotate_MissingKey_CreatesKey(t *testing.T) {
	opts := envrotator.DefaultOptions()
	opts.FailOnMissing = false
	result, err := envrotator.Rotate(base(), []envrotator.Rotation{
		{Key: "NEW_KEY", NewValue: "hello"},
	}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Vars["NEW_KEY"] != "hello" {
		t.Errorf("expected hello, got %q", result.Vars["NEW_KEY"])
	}
}

func TestRotate_FailOnMissing_ReturnsError(t *testing.T) {
	opts := envrotator.DefaultOptions()
	opts.FailOnMissing = true
	_, err := envrotator.Rotate(base(), []envrotator.Rotation{
		{Key: "DOES_NOT_EXIST", NewValue: "val"},
	}, opts)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestRotate_EmptyRotations_ReturnsUnchanged(t *testing.T) {
	v := base()
	result, err := envrotator.Rotate(v, nil, envrotator.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Log) != 0 {
		t.Errorf("expected empty log, got %d entries", len(result.Log))
	}
	if len(result.Vars) != len(v) {
		t.Errorf("expected %d vars, got %d", len(v), len(result.Vars))
	}
}

func TestRotate_EmptyKey_ReturnsError(t *testing.T) {
	_, err := envrotator.Rotate(base(), []envrotator.Rotation{
		{Key: "", NewValue: "val"},
	}, envrotator.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestRotate_MultipleKeys_AllRotated(t *testing.T) {
	result, err := envrotator.Rotate(base(), []envrotator.Rotation{
		{Key: "DB_PASSWORD", NewValue: "p1"},
		{Key: "API_KEY", NewValue: "k2"},
	}, envrotator.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Vars["DB_PASSWORD"] != "p1" {
		t.Errorf("DB_PASSWORD: expected p1, got %q", result.Vars["DB_PASSWORD"])
	}
	if result.Vars["API_KEY"] != "k2" {
		t.Errorf("API_KEY: expected k2, got %q", result.Vars["API_KEY"])
	}
	if len(result.Log) != 2 {
		t.Errorf("expected 2 log entries, got %d", len(result.Log))
	}
}
