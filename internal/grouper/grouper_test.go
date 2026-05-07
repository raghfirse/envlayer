package grouper_test

import (
	"testing"

	"github.com/yourusername/envlayer/internal/grouper"
)

func TestGroupBy_SplitsOnSeparator(t *testing.T) {
	vars := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "production",
	}

	groups := grouper.GroupBy(vars, "_")

	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	// Groups are sorted by name: APP, DB.
	if groups[0].Name != "APP" {
		t.Errorf("expected first group APP, got %q", groups[0].Name)
	}
	if groups[1].Name != "DB" {
		t.Errorf("expected second group DB, got %q", groups[1].Name)
	}
	if groups[1].Vars["HOST"] != "localhost" {
		t.Errorf("expected DB.HOST=localhost, got %q", groups[1].Vars["HOST"])
	}
	if groups[1].Vars["PORT"] != "5432" {
		t.Errorf("expected DB.PORT=5432, got %q", groups[1].Vars["PORT"])
	}
}

func TestGroupBy_RootGroupForKeysWithoutSeparator(t *testing.T) {
	vars := map[string]string{
		"PLAIN": "value",
		"DB_HOST": "localhost",
	}

	groups := grouper.GroupBy(vars, "_")

	var root *grouper.Group
	for i := range groups {
		if groups[i].Name == "" {
			root = &groups[i]
		}
	}
	if root == nil {
		t.Fatal("expected a root group for keys without separator")
	}
	if root.Vars["PLAIN"] != "value" {
		t.Errorf("expected root.PLAIN=value, got %q", root.Vars["PLAIN"])
	}
}

func TestGroupBy_EmptySeparatorReturnsSingleGroup(t *testing.T) {
	vars := map[string]string{"A": "1", "B": "2"}
	groups := grouper.GroupBy(vars, "")

	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Vars["A"] != "1" || groups[0].Vars["B"] != "2" {
		t.Error("root group does not contain all vars")
	}
}

func TestGroupBy_NestedSeparatorUsesFirstOccurrence(t *testing.T) {
	vars := map[string]string{"DB_HOST_PORT": "5432"}
	groups := grouper.GroupBy(vars, "_")

	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Name != "DB" {
		t.Errorf("expected group DB, got %q", groups[0].Name)
	}
	if groups[0].Vars["HOST_PORT"] != "5432" {
		t.Errorf("expected key HOST_PORT, got %v", groups[0].Vars)
	}
}

func TestFlatten_RoundTrip(t *testing.T) {
	original := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "staging",
		"PLAIN":   "yes",
	}

	groups := grouper.GroupBy(original, "_")
	result := grouper.Flatten(groups, "_")

	for k, v := range original {
		if result[k] != v {
			t.Errorf("key %q: expected %q, got %q", k, v, result[k])
		}
	}
	if len(result) != len(original) {
		t.Errorf("expected %d keys after flatten, got %d", len(original), len(result))
	}
}
