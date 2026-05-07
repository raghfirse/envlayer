package linter_test

import (
	"testing"

	"github.com/yourorg/envlayer/internal/linter"
)

func findingMessages(findings []linter.Finding) []string {
	msgs := make([]string, len(findings))
	for i, f := range findings {
		msgs[i] = f.String()
	}
	return msgs
}

func hasFinding(findings []linter.Finding, key string, sev linter.Severity) bool {
	for _, f := range findings {
		if f.Key == key && f.Severity == sev {
			return true
		}
	}
	return false
}

func TestLint_ValidVarsNoFindings(t *testing.T) {
	vars := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"PORT":         "8080",
	}
	findings := linter.Lint(vars)
	if len(findings) != 0 {
		t.Errorf("expected no findings, got: %v", findingMessages(findings))
	}
}

func TestLint_LowercaseKeyWarning(t *testing.T) {
	vars := map[string]string{"db_host": "localhost"}
	findings := linter.Lint(vars)
	if !hasFinding(findings, "db_host", linter.SeverityWarning) {
		t.Error("expected warning for lowercase key")
	}
}

func TestLint_EmptyValueWarning(t *testing.T) {
	vars := map[string]string{"SECRET_KEY": ""}
	findings := linter.Lint(vars)
	if !hasFinding(findings, "SECRET_KEY", linter.SeverityWarning) {
		t.Error("expected warning for empty value")
	}
}

func TestLint_EmptyOptionalKeyNoWarning(t *testing.T) {
	vars := map[string]string{"FEATURE_FLAG_OPTIONAL": ""}
	findings := linter.Lint(vars)
	for _, f := range findings {
		if f.Key == "FEATURE_FLAG_OPTIONAL" && f.Severity == linter.SeverityWarning {
			t.Error("should not warn on empty _OPTIONAL key")
		}
	}
}

func TestLint_UninterpolatedPlaceholderInfo(t *testing.T) {
	vars := map[string]string{"API_URL": "https://example.com/${PATH}"}
	findings := linter.Lint(vars)
	if !hasFinding(findings, "API_URL", linter.SeverityInfo) {
		t.Error("expected info finding for un-interpolated placeholder")
	}
}

func TestLint_FindingsAreSortedBySeverityThenKey(t *testing.T) {
	vars := map[string]string{
		"z_lower": "",  // warning (lowercase) + warning (empty)
		"A_EMPTY":  "", // warning (empty)
	}
	findings := linter.Lint(vars)
	// All warnings; A_EMPTY should appear before z_lower alphabetically within same severity.
	for i := 1; i < len(findings); i++ {
		prev, curr := findings[i-1], findings[i]
		if prev.Severity > curr.Severity {
			t.Errorf("findings not sorted by severity at index %d", i)
		}
	}
}

func TestFinding_String(t *testing.T) {
	f := linter.Finding{Key: "FOO", Message: "test message", Severity: linter.SeverityError}
	got := f.String()
	expected := "[error] FOO: test message"
	if got != expected {
		t.Errorf("String() = %q, want %q", got, expected)
	}
}
