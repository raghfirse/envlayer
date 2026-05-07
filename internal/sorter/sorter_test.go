package sorter_test

import (
	"testing"

	"github.com/user/envlayer/internal/sorter"
)

func TestByKey_Ascending(t *testing.T) {
	vars := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	entries := sorter.ByKey(vars, sorter.Ascending)
	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, e := range entries {
		if e.Key != expected[i] {
			t.Errorf("position %d: got %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestByKey_Descending(t *testing.T) {
	vars := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	entries := sorter.ByKey(vars, sorter.Descending)
	expected := []string{"ZEBRA", "MANGO", "APPLE"}
	for i, e := range entries {
		if e.Key != expected[i] {
			t.Errorf("position %d: got %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestByKey_EmptyMap(t *testing.T) {
	entries := sorter.ByKey(map[string]string{}, sorter.Ascending)
	if len(entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(entries))
	}
}

func TestByPriority_PriorityKeysFirst(t *testing.T) {
	vars := map[string]string{
		"APP_ENV": "prod",
		"DB_HOST": "localhost",
		"LOG_LEVEL": "info",
		"PORT": "8080",
	}
	priority := []string{"PORT", "APP_ENV"}
	entries := sorter.ByPriority(vars, priority)

	if entries[0].Key != "PORT" {
		t.Errorf("expected PORT first, got %q", entries[0].Key)
	}
	if entries[1].Key != "APP_ENV" {
		t.Errorf("expected APP_ENV second, got %q", entries[1].Key)
	}
	// remaining keys should be sorted alphabetically
	if entries[2].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST third, got %q", entries[2].Key)
	}
	if entries[3].Key != "LOG_LEVEL" {
		t.Errorf("expected LOG_LEVEL fourth, got %q", entries[3].Key)
	}
}

func TestByPriority_MissingPriorityKeySkipped(t *testing.T) {
	vars := map[string]string{"ALPHA": "1", "BETA": "2"}
	entries := sorter.ByPriority(vars, []string{"MISSING", "ALPHA"})
	if entries[0].Key != "ALPHA" {
		t.Errorf("expected ALPHA first, got %q", entries[0].Key)
	}
	if entries[1].Key != "BETA" {
		t.Errorf("expected BETA second, got %q", entries[1].Key)
	}
}

func TestToMap_RoundTrip(t *testing.T) {
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	entries := sorter.ByKey(vars, sorter.Ascending)
	result := sorter.ToMap(entries)
	for k, v := range vars {
		if result[k] != v {
			t.Errorf("key %q: got %q, want %q", k, result[k], v)
		}
	}
}
