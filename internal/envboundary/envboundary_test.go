package envboundary_test

import (
	"testing"

	"github.com/user/envlayer/internal/envboundary"
)

func makePolicy() envboundary.Policy {
	return envboundary.Policy{
		Rules: map[string]envboundary.Level{
			"SECRET_":    envboundary.Restricted,
			"INTERNAL_":  envboundary.Internal,
			"DB_PASSWORD": envboundary.Restricted,
		},
	}
}

func TestLevelFor_ExactKeyMatchTakesPrecedence(t *testing.T) {
	p := makePolicy()
	if got := p.LevelFor("DB_PASSWORD"); got != envboundary.Restricted {
		t.Fatalf("expected Restricted, got %d", got)
	}
}

func TestLevelFor_PrefixMatch(t *testing.T) {
	p := makePolicy()
	if got := p.LevelFor("SECRET_KEY"); got != envboundary.Restricted {
		t.Fatalf("expected Restricted, got %d", got)
	}
}

func TestLevelFor_DefaultsToPublic(t *testing.T) {
	p := makePolicy()
	if got := p.LevelFor("APP_NAME"); got != envboundary.Public {
		t.Fatalf("expected Public, got %d", got)
	}
}

func TestCheck_NoViolationsWhenAllAllowed(t *testing.T) {
	p := makePolicy()
	vars := map[string]string{
		"APP_NAME": "myapp",
		"APP_ENV":  "prod",
	}
	v := envboundary.Check(vars, p, envboundary.Public, envboundary.Internal, envboundary.Restricted)
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %d", len(v))
	}
}

func TestCheck_DetectsRestrictedKey(t *testing.T) {
	p := makePolicy()
	vars := map[string]string{
		"APP_NAME":   "myapp",
		"SECRET_KEY": "abc123",
	}
	v := envboundary.Check(vars, p, envboundary.Public)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Key != "SECRET_KEY" {
		t.Errorf("expected violation on SECRET_KEY, got %q", v[0].Key)
	}
}

func TestCheck_MultipleViolationsAreSorted(t *testing.T) {
	p := makePolicy()
	vars := map[string]string{
		"SECRET_B": "b",
		"SECRET_A": "a",
	}
	v := envboundary.Check(vars, p, envboundary.Public)
	if len(v) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(v))
	}
	if v[0].Key != "SECRET_A" || v[1].Key != "SECRET_B" {
		t.Errorf("violations not sorted: %v", v)
	}
}

func TestFilter_KeepsOnlyAllowedLevels(t *testing.T) {
	p := makePolicy()
	vars := map[string]string{
		"APP_NAME":    "myapp",
		"INTERNAL_ID": "42",
		"SECRET_KEY":  "s3cr3t",
	}
	out := envboundary.Filter(vars, p, envboundary.Public)
	if _, ok := out["APP_NAME"]; !ok {
		t.Error("expected APP_NAME in output")
	}
	if _, ok := out["INTERNAL_ID"]; ok {
		t.Error("INTERNAL_ID should be excluded")
	}
	if _, ok := out["SECRET_KEY"]; ok {
		t.Error("SECRET_KEY should be excluded")
	}
}

func TestFilter_EmptyVarsReturnsEmpty(t *testing.T) {
	p := makePolicy()
	out := envboundary.Filter(map[string]string{}, p, envboundary.Public)
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestViolation_ErrorMessage(t *testing.T) {
	v := envboundary.Violation{Key: "SECRET_X", Level: envboundary.Restricted, Reason: "level 2 not permitted in this context"}
	msg := v.Error()
	if msg == "" {
		t.Error("expected non-empty error message")
	}
}
