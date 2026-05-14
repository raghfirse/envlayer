package envfilter_test

import (
	"testing"

	"github.com/envlayer/envlayer/internal/envfilter"
)

func TestFilter_HasPrefix(t *testing.T) {
	vars := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "db",
	}
	got := envfilter.Filter(vars, envfilter.HasPrefix("APP_"))
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if _, ok := got["DB_HOST"]; ok {
		t.Error("DB_HOST should have been excluded")
	}
}

func TestFilter_HasSuffix(t *testing.T) {
	vars := map[string]string{
		"APP_HOST": "localhost",
		"DB_HOST":  "db",
		"APP_PORT": "8080",
	}
	got := envfilter.Filter(vars, envfilter.HasSuffix("_HOST"))
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if _, ok := got["APP_PORT"]; ok {
		t.Error("APP_PORT should have been excluded")
	}
}

func TestFilter_ValueNotEmpty(t *testing.T) {
	vars := map[string]string{
		"KEY_A": "value",
		"KEY_B": "",
		"KEY_C": "other",
	}
	got := envfilter.Filter(vars, envfilter.ValueNotEmpty())
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if _, ok := got["KEY_B"]; ok {
		t.Error("KEY_B should have been excluded")
	}
}

func TestFilter_KeyContains(t *testing.T) {
	vars := map[string]string{
		"AWS_SECRET_KEY": "secret",
		"AWS_ACCESS_KEY": "access",
		"APP_NAME":       "myapp",
	}
	got := envfilter.Filter(vars, envfilter.KeyContains("SECRET"))
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
}

func TestFilter_Not(t *testing.T) {
	vars := map[string]string{
		"APP_HOST": "localhost",
		"DB_HOST":  "db",
	}
	got := envfilter.Filter(vars, envfilter.Not(envfilter.HasPrefix("APP_")))
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if _, ok := got["DB_HOST"]; !ok {
		t.Error("DB_HOST should be present")
	}
}

func TestFilter_Any(t *testing.T) {
	vars := map[string]string{
		"APP_HOST": "localhost",
		"DB_HOST":  "db",
		"LOG_LEVEL": "info",
	}
	got := envfilter.Filter(vars, envfilter.Any(
		envfilter.HasPrefix("APP_"),
		envfilter.HasPrefix("DB_"),
	))
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if _, ok := got["LOG_LEVEL"]; ok {
		t.Error("LOG_LEVEL should have been excluded")
	}
}

func TestFilter_NoPredicates_ReturnsAll(t *testing.T) {
	vars := map[string]string{"A": "1", "B": "2"}
	got := envfilter.Filter(vars)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}
