package diff_test

import (
	"testing"

	"github.com/yourusername/envlayer/internal/diff"
)

func TestCompare_NoChanges(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	next := map[string]string{"A": "1", "B": "2"}
	entries := diff.Compare(base, next)
	if len(entries) != 0 {
		t.Fatalf("expected no entries, got %d", len(entries))
	}
}

func TestCompare_AddedKey(t *testing.T) {
	base := map[string]string{}
	next := map[string]string{"NEW": "hello"}
	entries := diff.Compare(base, next)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Type != diff.Added || e.Key != "NEW" || e.NewValue != "hello" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestCompare_RemovedKey(t *testing.T) {
	base := map[string]string{"OLD": "bye"}
	next := map[string]string{}
	entries := diff.Compare(base, next)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Type != diff.Removed || e.Key != "OLD" || e.OldValue != "bye" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestCompare_ChangedKey(t *testing.T) {
	base := map[string]string{"PORT": "3000"}
	next := map[string]string{"PORT": "8080"}
	entries := diff.Compare(base, next)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Type != diff.Changed || e.OldValue != "3000" || e.NewValue != "8080" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestCompare_SortedOutput(t *testing.T) {
	base := map[string]string{"Z": "z", "A": "a"}
	next := map[string]string{"Z": "z2", "A": "a"}
	entries := diff.Compare(base, next)
	if len(entries) != 1 || entries[0].Key != "Z" {
		t.Errorf("expected sorted single entry for Z, got %+v", entries)
	}
}

func TestCompare_MixedChanges(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2", "C": "3"}
	next := map[string]string{"A": "1", "B": "99", "D": "4"}
	entries := diff.Compare(base, next)
	// Expect: B changed, C removed, D added — 3 entries
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d: %+v", len(entries), entries)
	}
	types := map[string]diff.ChangeType{}
	for _, e := range entries {
		types[e.Key] = e.Type
	}
	if types["B"] != diff.Changed {
		t.Errorf("expected B to be Changed")
	}
	if types["C"] != diff.Removed {
		t.Errorf("expected C to be Removed")
	}
	if types["D"] != diff.Added {
		t.Errorf("expected D to be Added")
	}
}
