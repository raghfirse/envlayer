package envstats_test

import (
	"testing"

	"github.com/nicholasgasior/envlayer/internal/envstats"
)

func TestCompute_EmptyMap(t *testing.T) {
	s := envstats.Compute(map[string]string{})
	if s.TotalKeys != 0 {
		t.Fatalf("expected 0 keys, got %d", s.TotalKeys)
	}
	if s.PrefixGroups == nil {
		t.Fatal("PrefixGroups should be non-nil even for empty map")
	}
}

func TestCompute_TotalKeys(t *testing.T) {
	vars := map[string]string{"A": "1", "B": "2", "C": "3"}
	s := envstats.Compute(vars)
	if s.TotalKeys != 3 {
		t.Fatalf("expected 3, got %d", s.TotalKeys)
	}
}

func TestCompute_EmptyValues(t *testing.T) {
	vars := map[string]string{
		"PRESENT": "value",
		"EMPTY":   "",
		"ALSO":    "",
	}
	s := envstats.Compute(vars)
	if s.EmptyValues != 2 {
		t.Fatalf("expected 2 empty values, got %d", s.EmptyValues)
	}
}

func TestCompute_SensitiveKeys(t *testing.T) {
	vars := map[string]string{
		"DB_PASSWORD": "secret",
		"API_TOKEN":   "tok",
		"APP_NAME":    "myapp",
	}
	s := envstats.Compute(vars)
	if s.SensitiveKeys != 2 {
		t.Fatalf("expected 2 sensitive keys, got %d", s.SensitiveKeys)
	}
}

func TestCompute_AvgKeyLength(t *testing.T) {
	// keys: "AB" (2), "ABCD" (4) => avg 3.0
	vars := map[string]string{"AB": "x", "ABCD": "y"}
	s := envstats.Compute(vars)
	if s.AvgKeyLength != 3.0 {
		t.Fatalf("expected avg key length 3.0, got %f", s.AvgKeyLength)
	}
}

func TestCompute_PrefixGroups(t *testing.T) {
	vars := map[string]string{
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"APP_NAME": "envlayer",
		"NOPREFIX": "val",
	}
	s := envstats.Compute(vars)
	if s.PrefixGroups["DB"] != 2 {
		t.Fatalf("expected DB group count 2, got %d", s.PrefixGroups["DB"])
	}
	if s.PrefixGroups["APP"] != 1 {
		t.Fatalf("expected APP group count 1, got %d", s.PrefixGroups["APP"])
	}
	if s.PrefixGroups["NOPREFIX"] != 1 {
		t.Fatalf("expected NOPREFIX group count 1, got %d", s.PrefixGroups["NOPREFIX"])
	}
}

func TestTopPrefixes_SortedByCountDesc(t *testing.T) {
	vars := map[string]string{
		"DB_HOST":    "h",
		"DB_PORT":    "p",
		"DB_NAME":    "n",
		"APP_NAME":   "a",
		"APP_REGION": "r",
		"LOG_LEVEL":  "l",
	}
	s := envstats.Compute(vars)
	top := envstats.TopPrefixes(s)
	if len(top) == 0 {
		t.Fatal("expected non-empty top prefixes")
	}
	if top[0] != "DB" {
		t.Fatalf("expected DB to be top prefix, got %s", top[0])
	}
}

func TestTopPrefixes_EmptyStats(t *testing.T) {
	s := envstats.Compute(map[string]string{})
	top := envstats.TopPrefixes(s)
	if len(top) != 0 {
		t.Fatalf("expected empty top prefixes, got %v", top)
	}
}
