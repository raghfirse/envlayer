package differ_test

import (
	"testing"

	"github.com/user/envlayer/internal/differ"
)

func TestDiff_AddedKey(t *testing.T) {
	from := map[string]string{"A": "1"}
	to := map[string]string{"A": "1", "B": "2"}

	r := differ.Diff(from, to, "base", "prod")
	if len(r.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(r.Changes))
	}
	if r.Changes[0].Kind != differ.Added || r.Changes[0].Key != "B" {
		t.Errorf("unexpected change: %+v", r.Changes[0])
	}
}

func TestDiff_RemovedKey(t *testing.T) {
	from := map[string]string{"A": "1", "B": "2"}
	to := map[string]string{"A": "1"}

	r := differ.Diff(from, to, "base", "prod")
	if len(r.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(r.Changes))
	}
	if r.Changes[0].Kind != differ.Removed || r.Changes[0].Key != "B" {
		t.Errorf("unexpected change: %+v", r.Changes[0])
	}
}

func TestDiff_ChangedKey(t *testing.T) {
	from := map[string]string{"HOST": "localhost"}
	to := map[string]string{"HOST": "prod.example.com"}

	r := differ.Diff(from, to, "dev", "prod")
	if len(r.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(r.Changes))
	}
	c := r.Changes[0]
	if c.Kind != differ.Changed || c.OldValue != "localhost" || c.NewValue != "prod.example.com" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	r := differ.Diff(env, env, "x", "y")
	if len(r.Changes) != 0 {
		t.Errorf("expected no changes, got %d", len(r.Changes))
	}
}

func TestDiff_Labels(t *testing.T) {
	r := differ.Diff(map[string]string{}, map[string]string{}, "alpha", "beta")
	if r.From != "alpha" || r.To != "beta" {
		t.Errorf("wrong labels: %s %s", r.From, r.To)
	}
}

func TestSummary_Counts(t *testing.T) {
	from := map[string]string{"A": "1", "B": "old", "C": "3"}
	to := map[string]string{"B": "new", "C": "3", "D": "4"}

	r := differ.Diff(from, to, "dev", "prod")
	summary := differ.Summary(r)
	expected := "dev → prod: +1 -1 ~1"
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}

func TestDiff_SortedChanges(t *testing.T) {
	from := map[string]string{"Z": "1", "A": "old"}
	to := map[string]string{"Z": "1", "A": "new", "M": "mid"}

	r := differ.Diff(from, to, "a", "b")
	if len(r.Changes) < 2 {
		t.Fatalf("expected at least 2 changes")
	}
	if r.Changes[0].Key > r.Changes[1].Key {
		t.Errorf("changes not sorted: %s > %s", r.Changes[0].Key, r.Changes[1].Key)
	}
}
