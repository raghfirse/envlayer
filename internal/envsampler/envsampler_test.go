package envsampler_test

import (
	"testing"

	"github.com/nicholasgasior/envlayer/internal/envsampler"
)

var baseVars = map[string]string{
	"APP_HOST":     "localhost",
	"APP_PORT":     "8080",
	"DB_HOST":      "db.local",
	"DB_PASSWORD":  "secret",
	"LOG_LEVEL":    "info",
	"FEATURE_FLAG": "true",
}

func TestSample_ReturnsNEntries(t *testing.T) {
	opts := envsampler.DefaultOptions()
	opts.Seed = 42
	got, err := envsampler.Sample(baseVars, 3, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Errorf("expected 3 entries, got %d", len(got))
	}
}

func TestSample_NGreaterThanSize_ReturnsAll(t *testing.T) {
	got, err := envsampler.Sample(baseVars, 100, envsampler.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(baseVars) {
		t.Errorf("expected %d entries, got %d", len(baseVars), len(got))
	}
}

func TestSample_NegativeN_ReturnsError(t *testing.T) {
	_, err := envsampler.Sample(baseVars, -1, envsampler.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for negative n, got nil")
	}
}

func TestSample_WithPrefix_FiltersFirst(t *testing.T) {
	opts := envsampler.DefaultOptions()
	opts.Prefix = "DB_"
	opts.Seed = 1
	got, err := envsampler.Sample(baseVars, 10, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k := range got {
		if k[:3] != "DB_" {
			t.Errorf("unexpected key without DB_ prefix: %s", k)
		}
	}
	if len(got) != 2 {
		t.Errorf("expected 2 DB_ entries, got %d", len(got))
	}
}

func TestSample_DeterministicWithSameSeed(t *testing.T) {
	opts := envsampler.Options{Seed: 99}
	a, _ := envsampler.Sample(baseVars, 3, opts)
	b, _ := envsampler.Sample(baseVars, 3, opts)
	for k := range a {
		if _, ok := b[k]; !ok {
			t.Errorf("key %s present in first sample but not second", k)
		}
	}
}

func TestTopN_ReturnsLongestKeys(t *testing.T) {
	got, err := envsampler.TopN(baseVars, 2, envsampler.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 entries, got %d", len(got))
	}
	// FEATURE_FLAG (12) and DB_PASSWORD (11) are the two longest keys
	if _, ok := got["FEATURE_FLAG"]; !ok {
		t.Error("expected FEATURE_FLAG in top-2 longest keys")
	}
	if _, ok := got["DB_PASSWORD"]; !ok {
		t.Error("expected DB_PASSWORD in top-2 longest keys")
	}
}

func TestTopN_NegativeN_ReturnsError(t *testing.T) {
	_, err := envsampler.TopN(baseVars, -5, envsampler.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for negative n, got nil")
	}
}

func TestTopN_EmptyMap_ReturnsEmpty(t *testing.T) {
	got, err := envsampler.TopN(map[string]string{}, 3, envsampler.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %d entries", len(got))
	}
}
