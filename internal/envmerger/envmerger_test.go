package envmerger_test

import (
	"testing"

	"github.com/nicholasgasior/envlayer/internal/envmerger"
)

func TestMerge_LastWins_DefaultStrategy(t *testing.T) {
	a := map[string]string{"FOO": "base", "BAR": "bar"}
	b := map[string]string{"FOO": "override"}
	res, err := envmerger.Merge([]map[string]string{a, b}, envmerger.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["FOO"] != "override" {
		t.Errorf("expected FOO=override, got %q", res.Vars["FOO"])
	}
	if res.Vars["BAR"] != "bar" {
		t.Errorf("expected BAR=bar, got %q", res.Vars["BAR"])
	}
}

func TestMerge_FirstWins_DoesNotOverride(t *testing.T) {
	a := map[string]string{"FOO": "first"}
	b := map[string]string{"FOO": "second"}
	opts := envmerger.Options{Strategy: envmerger.FirstWins}
	res, err := envmerger.Merge([]map[string]string{a, b}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["FOO"] != "first" {
		t.Errorf("expected FOO=first, got %q", res.Vars["FOO"])
	}
}

func TestMerge_ErrorOnConflict_ReturnsError(t *testing.T) {
	a := map[string]string{"KEY": "v1"}
	b := map[string]string{"KEY": "v2"}
	opts := envmerger.Options{Strategy: envmerger.ErrorOnConflict}
	_, err := envmerger.Merge([]map[string]string{a, b}, opts)
	if err == nil {
		t.Fatal("expected error for conflicting key, got nil")
	}
}

func TestMerge_ErrorOnConflict_SameValueNoError(t *testing.T) {
	a := map[string]string{"KEY": "same"}
	b := map[string]string{"KEY": "same"}
	opts := envmerger.Options{Strategy: envmerger.ErrorOnConflict}
	_, err := envmerger.Merge([]map[string]string{a, b}, opts)
	if err != nil {
		t.Fatalf("unexpected error for identical values: %v", err)
	}
}

func TestMerge_ConflictsRecorded(t *testing.T) {
	a := map[string]string{"A": "1", "B": "x"}
	b := map[string]string{"A": "2"}
	res, err := envmerger.Merge([]map[string]string{a, b}, envmerger.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(res.Conflicts))
	}
	if res.Conflicts[0].Key != "A" {
		t.Errorf("expected conflict key A, got %q", res.Conflicts[0].Key)
	}
}

func TestMerge_EmptyLayers_ReturnsEmptyMap(t *testing.T) {
	res, err := envmerger.Merge([]map[string]string{}, envmerger.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 0 {
		t.Errorf("expected empty map, got %v", res.Vars)
	}
}

func TestMerge_NoConflicts_EmptyConflictSlice(t *testing.T) {
	a := map[string]string{"X": "1"}
	b := map[string]string{"Y": "2"}
	res, err := envmerger.Merge([]map[string]string{a, b}, envmerger.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
}
