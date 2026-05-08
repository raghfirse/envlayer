package envchain_test

import (
	"testing"

	"github.com/nicholasgasior/envlayer/internal/envchain"
)

func TestGet_HighestPriorityLayerWins(t *testing.T) {
	base := map[string]string{"APP_ENV": "base", "DB_HOST": "localhost"}
	override := map[string]string{"APP_ENV": "production"}
	c := envchain.New(override, base)

	got, err := c.Get("APP_ENV")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "production" {
		t.Errorf("expected production, got %s", got)
	}
}

func TestGet_FallsBackToLowerLayer(t *testing.T) {
	base := map[string]string{"DB_HOST": "localhost"}
	override := map[string]string{"APP_ENV": "production"}
	c := envchain.New(override, base)

	got, err := c.Get("DB_HOST")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "localhost" {
		t.Errorf("expected localhost, got %s", got)
	}
}

func TestGet_MissingKeyReturnsError(t *testing.T) {
	c := envchain.New(map[string]string{"A": "1"})
	_, err := c.Get("MISSING")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestGetOrDefault_ReturnsFallback(t *testing.T) {
	c := envchain.New(map[string]string{})
	got := c.GetOrDefault("NOT_HERE", "default_val")
	if got != "default_val" {
		t.Errorf("expected default_val, got %s", got)
	}
}

func TestResolve_MergesAllLayers(t *testing.T) {
	base := map[string]string{"A": "base_a", "B": "base_b"}
	top := map[string]string{"A": "top_a", "C": "top_c"}
	c := envchain.New(top, base)

	resolved := c.Resolve()
	if resolved["A"] != "top_a" {
		t.Errorf("A: expected top_a, got %s", resolved["A"])
	}
	if resolved["B"] != "base_b" {
		t.Errorf("B: expected base_b, got %s", resolved["B"])
	}
	if resolved["C"] != "top_c" {
		t.Errorf("C: expected top_c, got %s", resolved["C"])
	}
}

func TestPush_BecomesHighestPriority(t *testing.T) {
	base := map[string]string{"KEY": "old"}
	c := envchain.New(base)
	c.Push(map[string]string{"KEY": "new"})

	got, _ := c.Get("KEY")
	if got != "new" {
		t.Errorf("expected new, got %s", got)
	}
	if c.Len() != 2 {
		t.Errorf("expected 2 layers, got %d", c.Len())
	}
}

func TestNew_DoesNotMutateInput(t *testing.T) {
	original := map[string]string{"X": "1"}
	c := envchain.New(original)
	original["X"] = "mutated"

	got, _ := c.Get("X")
	if got != "1" {
		t.Errorf("chain was mutated by external map change, got %s", got)
	}
}
