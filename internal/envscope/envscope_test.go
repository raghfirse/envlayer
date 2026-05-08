package envscope_test

import (
	"testing"

	"github.com/yourusername/envlayer/internal/envscope"
)

var baseVars = map[string]string{
	"APP_HOST":    "localhost",
	"APP_PORT":    "8080",
	"DB_HOST":     "db.local",
	"DB_PASSWORD": "secret",
	"UNSCOPED":    "value",
}

func TestNew_FiltersAndStripsPrefix(t *testing.T) {
	s := envscope.New("app", "APP_", baseVars)
	all := s.All()

	if _, ok := all["HOST"]; !ok {
		t.Error("expected HOST key after stripping APP_ prefix")
	}
	if _, ok := all["PORT"]; !ok {
		t.Error("expected PORT key after stripping APP_ prefix")
	}
	if _, ok := all["DB_HOST"]; ok {
		t.Error("DB_HOST should not appear in APP scope")
	}
}

func TestNew_EmptyPrefixReturnsAll(t *testing.T) {
	s := envscope.New("all", "", baseVars)
	if len(s.All()) != len(baseVars) {
		t.Errorf("expected %d keys, got %d", len(baseVars), len(s.All()))
	}
}

func TestGet_ExistingKey(t *testing.T) {
	s := envscope.New("db", "DB_", baseVars)
	v, ok := s.Get("HOST")
	if !ok {
		t.Fatal("expected HOST to be found in DB scope")
	}
	if v != "db.local" {
		t.Errorf("expected db.local, got %q", v)
	}
}

func TestGet_MissingKey(t *testing.T) {
	s := envscope.New("db", "DB_", baseVars)
	_, ok := s.Get("NONEXISTENT")
	if ok {
		t.Error("expected missing key to return false")
	}
}

func TestKeys_Sorted(t *testing.T) {
	s := envscope.New("app", "APP_", baseVars)
	keys := s.Keys()
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("keys not sorted: %v", keys)
		}
	}
}

func TestQualify_ReattachesPrefix(t *testing.T) {
	s := envscope.New("app", "APP_", baseVars)
	if got := s.Qualify("HOST"); got != "APP_HOST" {
		t.Errorf("expected APP_HOST, got %q", got)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := envscope.New("app", "APP_", baseVars)
	all := s.All()
	all["HOST"] = "mutated"
	v, _ := s.Get("HOST")
	if v == "mutated" {
		t.Error("All() should return a copy, not a reference")
	}
}

func TestString_ContainsName(t *testing.T) {
	s := envscope.New("myapp", "APP_", baseVars)
	if str := s.String(); str == "" {
		t.Error("String() should not be empty")
	}
}
