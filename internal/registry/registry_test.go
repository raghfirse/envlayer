package registry

import (
	"testing"
)

func TestRegister_StoresEntry(t *testing.T) {
	r := New()
	r.Register("prod", map[string]string{"HOST": "example.com"}, "live")
	e, err := r.Get("prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Vars["HOST"] != "example.com" {
		t.Errorf("expected HOST=example.com, got %q", e.Vars["HOST"])
	}
}

func TestRegister_OverwritesExisting(t *testing.T) {
	r := New()
	r.Register("dev", map[string]string{"PORT": "3000"})
	r.Register("dev", map[string]string{"PORT": "4000"})
	e, _ := r.Get("dev")
	if e.Vars["PORT"] != "4000" {
		t.Errorf("expected PORT=4000, got %q", e.Vars["PORT"])
	}
}

func TestRegister_DoesNotMutateSource(t *testing.T) {
	r := New()
	src := map[string]string{"KEY": "val"}
	r.Register("test", src)
	src["KEY"] = "mutated"
	e, _ := r.Get("test")
	if e.Vars["KEY"] != "val" {
		t.Errorf("registry should store a copy, got %q", e.Vars["KEY"])
	}
}

func TestGet_MissingEntry_ReturnsError(t *testing.T) {
	r := New()
	_, err := r.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing entry")
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	r := New()
	r.Register("tmp", map[string]string{"A": "1"})
	if err := r.Remove("tmp"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err := r.Get("tmp")
	if err == nil {
		t.Fatal("expected entry to be removed")
	}
}

func TestRemove_MissingEntry_ReturnsError(t *testing.T) {
	r := New()
	if err := r.Remove("ghost"); err == nil {
		t.Fatal("expected error removing non-existent entry")
	}
}

func TestNames_ReturnsSorted(t *testing.T) {
	r := New()
	r.Register("zebra", map[string]string{})
	r.Register("alpha", map[string]string{})
	r.Register("mango", map[string]string{})
	names := r.Names()
	expected := []string{"alpha", "mango", "zebra"}
	for i, n := range expected {
		if names[i] != n {
			t.Errorf("index %d: expected %q, got %q", i, n, names[i])
		}
	}
}

func TestFindByTag_ReturnsMatchingEntries(t *testing.T) {
	r := New()
	r.Register("prod", map[string]string{}, "live", "stable")
	r.Register("staging", map[string]string{}, "live")
	r.Register("dev", map[string]string{}, "local")
	results := r.FindByTag("live")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Name != "prod" || results[1].Name != "staging" {
		t.Errorf("unexpected order: %v", results)
	}
}

func TestFindByTag_NoMatch_ReturnsEmpty(t *testing.T) {
	r := New()
	r.Register("dev", map[string]string{}, "local")
	results := r.FindByTag("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected empty result, got %d entries", len(results))
	}
}
