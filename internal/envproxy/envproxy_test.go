package envproxy_test

import (
	"testing"

	"github.com/your-org/envlayer/internal/envproxy"
)

func staticFallback(extra map[string]string) envproxy.FallbackFunc {
	return func(key string) (string, bool) {
		v, ok := extra[key]
		return v, ok
	}
}

func TestGet_PrimaryKeyReturned(t *testing.T) {
	p := envproxy.New(map[string]string{"FOO": "bar"}, nil)
	v, err := p.Get("FOO")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "bar" {
		t.Errorf("expected bar, got %q", v)
	}
}

func TestGet_FallsBackToFallback(t *testing.T) {
	p := envproxy.New(
		map[string]string{},
		staticFallback(map[string]string{"FROM_FALLBACK": "yes"}),
	)
	v, err := p.Get("FROM_FALLBACK")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "yes" {
		t.Errorf("expected yes, got %q", v)
	}
}

func TestGet_MissingKey_ReturnsError(t *testing.T) {
	p := envproxy.New(map[string]string{}, nil)
	_, err := p.Get("NOPE")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetOrDefault_ReturnsFallbackValue(t *testing.T) {
	p := envproxy.New(map[string]string{}, nil)
	v := p.GetOrDefault("MISSING", "default_val")
	if v != "default_val" {
		t.Errorf("expected default_val, got %q", v)
	}
}

func TestGetOrDefault_PrimaryTakesPrecedence(t *testing.T) {
	p := envproxy.New(map[string]string{"KEY": "primary"}, nil)
	v := p.GetOrDefault("KEY", "default_val")
	if v != "primary" {
		t.Errorf("expected primary, got %q", v)
	}
}

func TestHas_ExistingKey(t *testing.T) {
	p := envproxy.New(map[string]string{"A": "1"}, nil)
	if !p.Has("A") {
		t.Error("expected Has to return true")
	}
}

func TestHas_MissingKey(t *testing.T) {
	p := envproxy.New(map[string]string{}, nil)
	if p.Has("NOPE") {
		t.Error("expected Has to return false")
	}
}

func TestKeys_Sorted(t *testing.T) {
	p := envproxy.New(map[string]string{"Z": "1", "A": "2", "M": "3"}, nil)
	keys := p.Keys()
	expected := []string{"A", "M", "Z"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: expected %q got %q", i, expected[i], k)
		}
	}
}

func TestNew_DoesNotMutatePrimary(t *testing.T) {
	orig := map[string]string{"X": "1"}
	p := envproxy.New(orig, nil)
	orig["X"] = "mutated"
	v, _ := p.Get("X")
	if v != "1" {
		t.Errorf("proxy should not reflect mutation of original map, got %q", v)
	}
}

func TestResolve_ReturnsPrimarySnapshot(t *testing.T) {
	p := envproxy.New(map[string]string{"A": "1", "B": "2"}, nil)
	m := p.Resolve()
	if m["A"] != "1" || m["B"] != "2" {
		t.Errorf("unexpected resolve result: %v", m)
	}
}
