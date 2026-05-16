package envindexer_test

import (
	"testing"

	"github.com/your-org/envlayer/internal/envindexer"
)

func TestBuild_BasicIndex(t *testing.T) {
	vars := map[string]string{
		"HOST": "localhost",
		"DB_HOST": "localhost",
		"PORT": "5432",
	}
	idx := envindexer.Build(vars, envindexer.DefaultOptions())

	keys := envindexer.Lookup(idx, "localhost")
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys for 'localhost', got %d", len(keys))
	}
	if keys[0] != "DB_HOST" || keys[1] != "HOST" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestBuild_ExcludeEmpty(t *testing.T) {
	vars := map[string]string{
		"PRESENT": "value",
		"EMPTY":   "",
	}
	opts := envindexer.DefaultOptions()
	opts.ExcludeEmpty = true
	idx := envindexer.Build(vars, opts)

	if keys := envindexer.Lookup(idx, ""); keys != nil {
		t.Errorf("expected empty-value keys to be excluded, got %v", keys)
	}
	if keys := envindexer.Lookup(idx, "value"); len(keys) != 1 {
		t.Errorf("expected 1 key for 'value', got %v", keys)
	}
}

func TestBuild_IncludeEmpty_WhenDisabled(t *testing.T) {
	vars := map[string]string{
		"EMPTY_A": "",
		"EMPTY_B": "",
	}
	opts := envindexer.DefaultOptions()
	opts.ExcludeEmpty = false
	idx := envindexer.Build(vars, opts)

	keys := envindexer.Lookup(idx, "")
	if len(keys) != 2 {
		t.Errorf("expected 2 keys for empty value, got %v", keys)
	}
}

func TestBuild_CaseFoldValues(t *testing.T) {
	vars := map[string]string{
		"ENV": "Production",
		"APP_ENV": "production",
	}
	opts := envindexer.DefaultOptions()
	opts.CaseFoldValues = true
	idx := envindexer.Build(vars, opts)

	keys := envindexer.Lookup(idx, "production")
	if len(keys) != 2 {
		t.Errorf("expected 2 keys after case-fold, got %v", keys)
	}
}

func TestLookup_MissingValue_ReturnsNil(t *testing.T) {
	idx := envindexer.Build(map[string]string{"K": "v"}, envindexer.DefaultOptions())
	if keys := envindexer.Lookup(idx, "missing"); keys != nil {
		t.Errorf("expected nil for missing value, got %v", keys)
	}
}

func TestValues_ReturnsSorted(t *testing.T) {
	vars := map[string]string{
		"C": "zebra",
		"A": "apple",
		"B": "mango",
	}
	idx := envindexer.Build(vars, envindexer.DefaultOptions())
	vals := envindexer.Values(idx)

	for i := 1; i < len(vals); i++ {
		if vals[i-1] > vals[i] {
			t.Errorf("values not sorted: %v", vals)
		}
	}
}

func TestLookup_DoesNotMutateIndex(t *testing.T) {
	vars := map[string]string{"X": "shared", "Y": "shared"}
	idx := envindexer.Build(vars, envindexer.DefaultOptions())

	first := envindexer.Lookup(idx, "shared")
	first[0] = "MUTATED"

	second := envindexer.Lookup(idx, "shared")
	for _, k := range second {
		if k == "MUTATED" {
			t.Error("Lookup returned a slice that shares backing array with index")
		}
	}
}
