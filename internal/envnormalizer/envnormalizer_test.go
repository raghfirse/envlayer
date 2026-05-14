package envnormalizer_test

import (
	"testing"

	"github.com/yourusername/envlayer/internal/envnormalizer"
)

func TestNormalize_TrimValues(t *testing.T) {
	src := map[string]string{"KEY": "  hello  ", "B": "\tworld\t"}
	opts := envnormalizer.DefaultOptions()
	opts.TrimValues = true
	out := envnormalizer.Normalize(src, opts)
	if out["KEY"] != "hello" {
		t.Errorf("expected 'hello', got %q", out["KEY"])
	}
	if out["B"] != "world" {
		t.Errorf("expected 'world', got %q", out["B"])
	}
}

func TestNormalize_UppercaseKeys(t *testing.T) {
	src := map[string]string{"db_host": "localhost", "db_port": "5432"}
	opts := envnormalizer.DefaultOptions()
	opts.UppercaseKeys = true
	out := envnormalizer.Normalize(src, opts)
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected key DB_HOST to exist")
	}
	if _, ok := out["DB_PORT"]; !ok {
		t.Error("expected key DB_PORT to exist")
	}
	if _, ok := out["db_host"]; ok {
		t.Error("original lowercase key should not exist")
	}
}

func TestNormalize_LowercaseKeys(t *testing.T) {
	src := map[string]string{"APP_NAME": "envlayer"}
	opts := envnormalizer.DefaultOptions()
	opts.LowercaseKeys = true
	out := envnormalizer.Normalize(src, opts)
	if out["app_name"] != "envlayer" {
		t.Errorf("expected app_name=envlayer, got %v", out)
	}
}

func TestNormalize_UppercaseTakesPrecedenceOverLowercase(t *testing.T) {
	src := map[string]string{"mixed_Key": "val"}
	opts := envnormalizer.Options{UppercaseKeys: true, LowercaseKeys: true}
	out := envnormalizer.Normalize(src, opts)
	if _, ok := out["MIXED_KEY"]; !ok {
		t.Error("expected MIXED_KEY when both uppercase and lowercase are set")
	}
}

func TestNormalize_RemoveEmpty(t *testing.T) {
	src := map[string]string{"PRESENT": "yes", "EMPTY": "", "SPACES": "   "}
	opts := envnormalizer.Options{TrimValues: true, RemoveEmpty: true}
	out := envnormalizer.Normalize(src, opts)
	if _, ok := out["EMPTY"]; ok {
		t.Error("EMPTY key should have been removed")
	}
	if _, ok := out["SPACES"]; ok {
		t.Error("SPACES key should have been removed after trim")
	}
	if out["PRESENT"] != "yes" {
		t.Errorf("PRESENT should remain, got %q", out["PRESENT"])
	}
}

func TestNormalize_CollapseWhitespace(t *testing.T) {
	src := map[string]string{"MSG": "hello   world\t\there"}
	opts := envnormalizer.Options{CollapseWhitespace: true}
	out := envnormalizer.Normalize(src, opts)
	if out["MSG"] != "hello world here" {
		t.Errorf("unexpected value: %q", out["MSG"])
	}
}

func TestNormalize_DoesNotMutateSource(t *testing.T) {
	src := map[string]string{"KEY": "  value  "}
	opts := envnormalizer.DefaultOptions()
	envnormalizer.Normalize(src, opts)
	if src["KEY"] != "  value  " {
		t.Error("source map was mutated")
	}
}

func TestNormalize_EmptyMap(t *testing.T) {
	out := envnormalizer.Normalize(map[string]string{}, envnormalizer.DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
