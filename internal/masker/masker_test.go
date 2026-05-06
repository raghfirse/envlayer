package masker_test

import (
	"testing"

	"github.com/yourusername/envlayer/internal/masker"
)

func TestMask_SensitiveKeysAreReplaced(t *testing.T) {
	vars := map[string]string{
		"DB_PASSWORD": "supersecret",
		"APP_NAME":    "myapp",
	}
	result := masker.Mask(vars, masker.Options{})

	if result["DB_PASSWORD"] != "****" {
		t.Errorf("expected DB_PASSWORD to be masked, got %q", result["DB_PASSWORD"])
	}
	if result["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME to be unchanged, got %q", result["APP_NAME"])
	}
}

func TestMask_CustomMaskChar(t *testing.T) {
	vars := map[string]string{"API_KEY": "abc123"}
	result := masker.Mask(vars, masker.Options{MaskChar: "[REDACTED]"})

	if result["API_KEY"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", result["API_KEY"])
	}
}

func TestMask_CustomSensitiveKeys(t *testing.T) {
	vars := map[string]string{
		"STRIPE_TOKEN": "tok_live_xyz",
		"DATABASE_URL": "postgres://localhost/db",
		"MY_INTERNAL":  "value",
	}
	result := masker.Mask(vars, masker.Options{
		SensitiveKeys: []string{"INTERNAL"},
	})

	if result["MY_INTERNAL"] != "****" {
		t.Errorf("expected MY_INTERNAL to be masked")
	}
	// STRIPE_TOKEN not in custom list
	if result["STRIPE_TOKEN"] != "tok_live_xyz" {
		t.Errorf("expected STRIPE_TOKEN to be unmasked")
	}
}

func TestMask_CaseInsensitiveMatching(t *testing.T) {
	vars := map[string]string{"db_password": "secret"}
	result := masker.Mask(vars, masker.Options{})

	if result["db_password"] != "****" {
		t.Errorf("expected lowercase key to be masked")
	}
}

func TestMask_EmptyVars(t *testing.T) {
	result := masker.Mask(map[string]string{}, masker.Options{})
	if len(result) != 0 {
		t.Errorf("expected empty result for empty input")
	}
}

func TestIsSensitive(t *testing.T) {
	cases := []struct {
		key      string
		expected bool
	}{
		{"DB_PASSWORD", true},
		{"API_KEY", true},
		{"APP_NAME", false},
		{"PORT", false},
		{"AUTH_TOKEN", true},
	}
	for _, tc := range cases {
		got := masker.IsSensitive(tc.key, masker.DefaultSensitiveKeys)
		if got != tc.expected {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.expected)
		}
	}
}
