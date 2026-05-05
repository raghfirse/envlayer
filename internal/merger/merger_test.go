package merger_test

import (
	"reflect"
	"testing"

	"github.com/envlayer/envlayer/internal/merger"
)

func TestMerge_SingleLayer(t *testing.T) {
	layers := []map[string]string{
		{"APP_ENV": "development", "PORT": "8080"},
	}
	got := merger.Merge(layers)
	want := map[string]string{"APP_ENV": "development", "PORT": "8080"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Merge() = %v, want %v", got, want)
	}
}

func TestMerge_OverridesInOrder(t *testing.T) {
	layers := []map[string]string{
		{"APP_ENV": "development", "PORT": "8080", "DB_HOST": "localhost"},
		{"APP_ENV": "staging", "PORT": "9090"},
	}
	got := merger.Merge(layers)
	want := map[string]string{
		"APP_ENV": "staging",
		"PORT":    "9090",
		"DB_HOST": "localhost",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Merge() = %v, want %v", got, want)
	}
}

func TestMerge_EmptyLayers(t *testing.T) {
	got := merger.Merge([]map[string]string{})
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestMergeWithBase_DoesNotMutateBase(t *testing.T) {
	base := map[string]string{"KEY": "base", "ONLY_BASE": "yes"}
	override := map[string]string{"KEY": "overridden"}
	got := merger.MergeWithBase(base, override)

	if base["KEY"] != "base" {
		t.Error("MergeWithBase mutated the base map")
	}
	if got["KEY"] != "overridden" {
		t.Errorf("expected KEY=overridden, got %s", got["KEY"])
	}
	if got["ONLY_BASE"] != "yes" {
		t.Errorf("expected ONLY_BASE=yes, got %s", got["ONLY_BASE"])
	}
}

func TestKeys_Sorted(t *testing.T) {
	env := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	got := merger.Keys(env)
	want := []string{"APPLE", "MANGO", "ZEBRA"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Keys() = %v, want %v", got, want)
	}
}
