package envdefaults_test

import (
	"testing"

	"github.com/yourorg/envlayer/internal/envdefaults"
)

func TestApply_MissingKeyIsSet(t *testing.T) {
	target := map[string]string{"A": "1"}
	defaults := map[string]string{"A": "99", "B": "2"}

	out, report, err := envdefaults.Apply(target, defaults, envdefaults.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" {
		t.Errorf("expected A=1, got %q", out["A"])
	}
	if out["B"] != "2" {
		t.Errorf("expected B=2, got %q", out["B"])
	}
	if len(report.Set) != 1 || report.Set[0] != "B" {
		t.Errorf("expected Set=[B], got %v", report.Set)
	}
	if len(report.Skipped) != 1 || report.Skipped[0] != "A" {
		t.Errorf("expected Skipped=[A], got %v", report.Skipped)
	}
}

func TestApply_EmptyValueOverwritten(t *testing.T) {
	target := map[string]string{"HOST": ""}
	defaults := map[string]string{"HOST": "localhost"}

	out, report, err := envdefaults.Apply(target, defaults, envdefaults.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", out["HOST"])
	}
	if len(report.Set) != 1 || report.Set[0] != "HOST" {
		t.Errorf("expected Set=[HOST], got %v", report.Set)
	}
}

func TestApply_EmptyValueNotOverwrittenWhenDisabled(t *testing.T) {
	opts := envdefaults.Options{OverwriteEmpty: false, FailOnConflict: false}
	target := map[string]string{"PORT": ""}
	defaults := map[string]string{"PORT": "8080"}

	out, _, err := envdefaults.Apply(target, defaults, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PORT"] != "" {
		t.Errorf("expected PORT empty, got %q", out["PORT"])
	}
}

func TestApply_FailOnConflict_ReturnsError(t *testing.T) {
	opts := envdefaults.Options{OverwriteEmpty: true, FailOnConflict: true}
	target := map[string]string{"KEY": "original"}
	defaults := map[string]string{"KEY": "default"}

	_, _, err := envdefaults.Apply(target, defaults, opts)
	if err == nil {
		t.Fatal("expected error for conflict, got nil")
	}
}

func TestApply_FailOnConflict_SameValueNoError(t *testing.T) {
	opts := envdefaults.Options{OverwriteEmpty: true, FailOnConflict: true}
	target := map[string]string{"KEY": "same"}
	defaults := map[string]string{"KEY": "same"}

	_, _, err := envdefaults.Apply(target, defaults, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_DoesNotMutateTarget(t *testing.T) {
	target := map[string]string{"X": "1"}
	defaults := map[string]string{"Y": "2"}

	_, _, _ = envdefaults.Apply(target, defaults, envdefaults.DefaultOptions())
	if _, ok := target["Y"]; ok {
		t.Error("Apply must not mutate the target map")
	}
}

func TestApply_EmptyDefaults_ReturnsClone(t *testing.T) {
	target := map[string]string{"A": "1", "B": "2"}

	out, report, err := envdefaults.Apply(target, map[string]string{}, envdefaults.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if len(report.Set) != 0 || len(report.Skipped) != 0 {
		t.Errorf("expected empty report, got %+v", report)
	}
}
