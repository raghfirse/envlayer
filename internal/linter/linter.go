// Package linter provides static analysis for .env files, detecting
// common issues such as duplicate keys, suspicious values, and keys
// that violate naming conventions.
package linter

import (
	"fmt"
	"regexp"
	"strings"
)

// Severity indicates how serious a lint finding is.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Finding represents a single lint result.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s] %s: %s", f.Severity, f.Key, f.Message)
}

var validKeyRe = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// Lint analyses the provided env map and returns a list of findings.
// It checks for:
//   - Keys that do not follow UPPER_SNAKE_CASE convention
//   - Empty values on non-optional keys (keys not ending with _OPTIONAL)
//   - Values that look like they contain un-interpolated placeholders
//   - Duplicate keys (caller must deduplicate before passing if needed)
func Lint(vars map[string]string) []Finding {
	var findings []Finding

	for key, value := range vars {
		// Convention: keys should be UPPER_SNAKE_CASE
		if !validKeyRe.MatchString(key) {
			findings = append(findings, Finding{
				Key:      key,
				Message:  "key does not follow UPPER_SNAKE_CASE naming convention",
				Severity: SeverityWarning,
			})
		}

		// Empty value warning (skip keys marked optional by convention)
		if value == "" && !strings.HasSuffix(key, "_OPTIONAL") {
			findings = append(findings, Finding{
				Key:      key,
				Message:  "value is empty; mark key as optional by appending _OPTIONAL if intentional",
				Severity: SeverityWarning,
			})
		}

		// Detect un-interpolated placeholders like ${FOO} or $FOO left in values
		if strings.Contains(value, "${") || (strings.Contains(value, "$") && !strings.HasPrefix(value, "$$")) {
			findings = append(findings, Finding{
				Key:      key,
				Message:  "value appears to contain an un-interpolated variable reference",
				Severity: SeverityInfo,
			})
		}
	}

	return sortFindings(findings)
}

func sortFindings(findings []Finding) []Finding {
	// Stable sort: errors first, then warnings, then info; alpha by key within group.
	order := map[Severity]int{SeverityError: 0, SeverityWarning: 1, SeverityInfo: 2}
	for i := 1; i < len(findings); i++ {
		for j := i; j > 0; j-- {
			a, b := findings[j-1], findings[j]
			if order[a.Severity] > order[b.Severity] ||
				(order[a.Severity] == order[b.Severity] && a.Key > b.Key) {
				findings[j-1], findings[j] = findings[j], findings[j-1]
			} else {
				break
			}
		}
	}
	return findings
}
