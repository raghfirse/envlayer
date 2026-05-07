package transformer_test

import (
	"testing"

	"github.com/yourorg/envlayer/internal/transformer"
)

func TestTransform_AddPrefix(t *testing.T) {
	vars := map[string]string{"HOST": "localhost", "PORT": "5432"}
	out := transformer.Transform(vars, transformer.Options{AddPrefix: "APP_"})
	if out["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %q", out["APP_HOST"])
	}
	if out["APP_PORT"] != "5432" {
		t.Errorf("expected APP_PORT=5432, got %q", out["APP_PORT"])
	}
	if _, ok := out["HOST"]; ok {
		t.Error("original key HOST should not exist after prefix add")
	}
}

func TestTransform_StripPrefix(t *testing.T) {
	vars := map[string]string{"APP_HOST": "localhost", "APP_PORT": "5432"}
	out := transformer.Transform(vars, transformer.Options{StripPrefix: "APP_"})
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", out["HOST"])
	}
	if out["PORT"] != "5432" {
		t.Errorf("expected PORT=5432, got %q", out["PORT"])
	}
}

func TestTransform_UppercaseKeys(t *testing.T) {
	vars := map[string]string{"db_host": "localhost"}
	out := transformer.Transform(vars, transformer.Options{UppercaseKeys: true})
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
}

func TestTransform_LowercaseKeys(t *testing.T) {
	vars := map[string]string{"DB_HOST": "localhost"}
	out := transformer.Transform(vars, transformer.Options{LowercaseKeys: true})
	if out["db_host"] != "localhost" {
		t.Errorf("expected db_host=localhost, got %q", out["db_host"])
	}
}

func TestTransform_TrimValues(t *testing.T) {
	vars := map[string]string{"KEY": "  value  "}
	out := transformer.Transform(vars, transformer.Options{TrimValues: true})
	if out["KEY"] != "value" {
		t.Errorf("expected trimmed value, got %q", out["KEY"])
	}
}

func TestTransform_DoesNotMutateOriginal(t *testing.T) {
	vars := map[string]string{"HOST": "localhost"}
	transformer.Transform(vars, transformer.Options{AddPrefix: "X_", UppercaseKeys: true})
	if _, ok := vars["HOST"]; !ok {
		t.Error("original map should not be mutated")
	}
}

func TestTransform_NoOptions_ReturnsShallowCopy(t *testing.T) {
	vars := map[string]string{"A": "1", "B": "2"}
	out := transformer.Transform(vars, transformer.Options{})
	if len(out) != len(vars) {
		t.Errorf("expected %d keys, got %d", len(vars), len(out))
	}
	for k, v := range vars {
		if out[k] != v {
			t.Errorf("key %s: expected %q got %q", k, v, out[k])
		}
	}
}
