package envcompare_test

import (
	"testing"

	"github.com/envlayer/envlayer/internal/envcompare"
)

func TestCompare_IdenticalMaps(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1", "B": "2"}
	r := envcompare.Compare(left, right)
	if len(r.OnlyInLeft) != 0 || len(r.OnlyInRight) != 0 || len(r.Changed) != 0 {
		t.Errorf("expected no differences, got %+v", r)
	}
	if len(r.Identical) != 2 {
		t.Errorf("expected 2 identical keys, got %d", len(r.Identical))
	}
}

func TestCompare_AddedKeys(t *testing.T) {
	left := map[string]string{"A": "1"}
	right := map[string]string{"A": "1", "B": "2"}
	r := envcompare.Compare(left, right)
	if len(r.OnlyInRight) != 1 || r.OnlyInRight[0] != "B" {
		t.Errorf("expected B in OnlyInRight, got %v", r.OnlyInRight)
	}
}

func TestCompare_RemovedKeys(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1"}
	r := envcompare.Compare(left, right)
	if len(r.OnlyInLeft) != 1 || r.OnlyInLeft[0] != "B" {
		t.Errorf("expected B in OnlyInLeft, got %v", r.OnlyInLeft)
	}
}

func TestCompare_ChangedValues(t *testing.T) {
	left := map[string]string{"A": "old", "B": "same"}
	right := map[string]string{"A": "new", "B": "same"}
	r := envcompare.Compare(left, right)
	pair, ok := r.Changed["A"]
	if !ok {
		t.Fatal("expected A in Changed")
	}
	if pair[0] != "old" || pair[1] != "new" {
		t.Errorf("expected [old new], got %v", pair)
	}
	if len(r.Identical) != 1 || r.Identical[0] != "B" {
		t.Errorf("expected B identical, got %v", r.Identical)
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	r := envcompare.Compare(map[string]string{}, map[string]string{})
	if len(r.OnlyInLeft)+len(r.OnlyInRight)+len(r.Changed)+len(r.Identical) != 0 {
		t.Errorf("expected empty result for empty maps")
	}
}

func TestEqual_SameMaps(t *testing.T) {
	if !envcompare.Equal(map[string]string{"X": "1"}, map[string]string{"X": "1"}) {
		t.Error("expected Equal to return true")
	}
}

func TestEqual_DifferentValues(t *testing.T) {
	if envcompare.Equal(map[string]string{"X": "1"}, map[string]string{"X": "2"}) {
		t.Error("expected Equal to return false")
	}
}

func TestEqual_DifferentLengths(t *testing.T) {
	if envcompare.Equal(map[string]string{"X": "1"}, map[string]string{"X": "1", "Y": "2"}) {
		t.Error("expected Equal to return false for different lengths")
	}
}

func TestSummary_NoDifferences(t *testing.T) {
	r := envcompare.Compare(map[string]string{"A": "1"}, map[string]string{"A": "1"})
	s := envcompare.Summary(r)
	if s != "no differences found" {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestSummary_WithChanges(t *testing.T) {
	left := map[string]string{"A": "old", "B": "gone"}
	right := map[string]string{"A": "new", "C": "added"}
	r := envcompare.Compare(left, right)
	s := envcompare.Summary(r)
	if s == "" || s == "no differences found" {
		t.Errorf("expected non-empty diff summary, got %q", s)
	}
}
