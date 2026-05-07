package envmap_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlayer/internal/envmap"
)

func TestFromMap_SortedKeys(t *testing.T) {
	m := map[string]string{"ZEBRA": "z", "ALPHA": "a", "MIDDLE": "m"}
	entries := envmap.FromMap(m)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Key != "ALPHA" || entries[1].Key != "MIDDLE" || entries[2].Key != "ZEBRA" {
		t.Errorf("unexpected order: %v", entries)
	}
}

func TestToMap_RoundTrip(t *testing.T) {
	orig := map[string]string{"A": "1", "B": "2"}
	entries := envmap.FromMap(orig)
	result := envmap.ToMap(entries)
	for k, v := range orig {
		if result[k] != v {
			t.Errorf("key %s: expected %q, got %q", k, v, result[k])
		}
	}
}

func TestToMap_LaterEntryWins(t *testing.T) {
	entries := []envmap.Entry{
		{Key: "X", Value: "first"},
		{Key: "X", Value: "second"},
	}
	m := envmap.ToMap(entries)
	if m["X"] != "second" {
		t.Errorf("expected 'second', got %q", m["X"])
	}
}

func TestFilter_SelectsMatchingKeys(t *testing.T) {
	entries := []envmap.Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "DB_HOST", Value: "db"},
		{Key: "APP_PORT", Value: "8080"},
	}
	filtered := envmap.Filter(entries, func(k string) bool {
		return strings.HasPrefix(k, "APP_")
	})
	if len(filtered) != 2 {
		t.Fatalf("expected 2, got %d", len(filtered))
	}
	for _, e := range filtered {
		if !strings.HasPrefix(e.Key, "APP_") {
			t.Errorf("unexpected key %q in filtered result", e.Key)
		}
	}
}

func TestMapValues_TransformsValues(t *testing.T) {
	entries := []envmap.Entry{
		{Key: "NAME", Value: "world"},
		{Key: "GREETING", Value: "hello"},
	}
	upper := envmap.MapValues(entries, func(_, v string) string {
		return strings.ToUpper(v)
	})
	if upper[0].Value != "WORLD" {
		t.Errorf("expected WORLD, got %q", upper[0].Value)
	}
	if upper[1].Value != "HELLO" {
		t.Errorf("expected HELLO, got %q", upper[1].Value)
	}
}

func TestFilter_EmptyInput(t *testing.T) {
	result := envmap.Filter(nil, func(k string) bool { return true })
	if len(result) != 0 {
		t.Errorf("expected empty, got %d entries", len(result))
	}
}
