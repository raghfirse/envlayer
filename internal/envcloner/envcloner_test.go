package envcloner_test

import (
	"testing"

	"github.com/nicholasgasior/envlayer/internal/envcloner"
)

func TestClone_DeepCopy(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	cloned := envcloner.Clone(src, envcloner.DefaultOptions())
	cloned["A"] = "mutated"
	if src["A"] != "1" {
		t.Errorf("Clone mutated original: got %q", src["A"])
	}
	if cloned["B"] != "2" {
		t.Errorf("expected B=2, got %q", cloned["B"])
	}
}

func TestClone_KeyPrefix(t *testing.T) {
	src := map[string]string{"HOST": "localhost", "PORT": "5432"}
	opts := envcloner.Options{KeyPrefix: "DB_"}
	out := envcloner.Clone(src, opts)
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
	if _, ok := out["HOST"]; ok {
		t.Error("original key HOST should not appear in output")
	}
}

func TestClone_StripPrefix(t *testing.T) {
	src := map[string]string{"APP_NAME": "envlayer", "APP_ENV": "prod"}
	opts := envcloner.Options{StripPrefix: "APP_"}
	out := envcloner.Clone(src, opts)
	if out["NAME"] != "envlayer" {
		t.Errorf("expected NAME=envlayer, got %q", out["NAME"])
	}
	if out["ENV"] != "prod" {
		t.Errorf("expected ENV=prod, got %q", out["ENV"])
	}
}

func TestClone_UppercaseKeys(t *testing.T) {
	src := map[string]string{"host": "localhost", "port": "80"}
	opts := envcloner.Options{UppercaseKeys: true}
	out := envcloner.Clone(src, opts)
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", out["HOST"])
	}
	if out["PORT"] != "80" {
		t.Errorf("expected PORT=80, got %q", out["PORT"])
	}
}

func TestClone_OmitEmpty(t *testing.T) {
	src := map[string]string{"KEY": "value", "EMPTY": ""}
	opts := envcloner.Options{OmitEmpty: true}
	out := envcloner.Clone(src, opts)
	if _, ok := out["EMPTY"]; ok {
		t.Error("expected EMPTY key to be omitted")
	}
	if out["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", out["KEY"])
	}
}

func TestCloneKeys_SelectsSubset(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	out := envcloner.CloneKeys(src, []string{"A", "C"})
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if out["A"] != "1" || out["C"] != "3" {
		t.Errorf("unexpected values: %v", out)
	}
}

func TestCloneKeys_MissingKeySkipped(t *testing.T) {
	src := map[string]string{"A": "1"}
	out := envcloner.CloneKeys(src, []string{"A", "MISSING"})
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
}

func TestCloneExclude_OmitsKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	out := envcloner.CloneExclude(src, []string{"B"})
	if _, ok := out["B"]; ok {
		t.Error("expected B to be excluded")
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}
