// Package envscorer computes a quality score for an environment variable map
// based on configurable heuristics such as key naming conventions, value
// completeness, and sensitive key coverage.
package envscorer

import (
	"strings"
)

// Result holds the overall score and per-category breakdown.
type Result struct {
	Total       int            // 0–100
	Categories  map[string]int // category name → score contribution
	Deductions  []string       // human-readable reasons for point loss
}

// Options controls which heuristics are applied.
type Options struct {
	PenalizeEmptyValues  bool // deduct points for keys with empty values
	PenalizeLowercase    bool // deduct points for keys that are not UPPER_SNAKE_CASE
	PenalizeShortKeys    bool // deduct points for keys shorter than MinKeyLen
	MinKeyLen            int  // minimum acceptable key length (default 3)
	SensitivePatterns    []string // keys matching these substrings should be non-empty
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		PenalizeEmptyValues: true,
		PenalizeLowercase:   true,
		PenalizeShortKeys:   true,
		MinKeyLen:           3,
		SensitivePatterns:   []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PASS"},
	}
}

// Score evaluates vars and returns a Result.
func Score(vars map[string]string, opts Options) Result {
	if opts.MinKeyLen == 0 {
		opts.MinKeyLen = 3
	}

	total := 100
	var deductions []string
	categories := map[string]int{
		"naming":     25,
		"completeness": 50,
		"security":   25,
	}

	if len(vars) == 0 {
		return Result{Total: 0, Categories: map[string]int{"naming": 0, "completeness": 0, "security": 0}, Deductions: []string{"no variables defined"}}
	}

	// Naming checks
	namingPenalty := 0
	for k := range vars {
		if opts.PenalizeLowercase && k != strings.ToUpper(k) {
			namingPenalty++
			deductions = append(deductions, "key not uppercase: "+k)
		}
		if opts.PenalizeShortKeys && len(k) < opts.MinKeyLen {
			namingPenalty++
			deductions = append(deductions, "key too short: "+k)
		}
	}
	namingScore := 25
	if len(vars) > 0 && namingPenalty > 0 {
		cut := (namingPenalty * 25) / len(vars)
		if cut > 25 {
			cut = 25
		}
		namingScore -= cut
		total -= cut
	}
	categories["naming"] = namingScore

	// Completeness checks
	emptyCount := 0
	if opts.PenalizeEmptyValues {
		for k, v := range vars {
			if strings.TrimSpace(v) == "" {
				emptyCount++
				deductions = append(deductions, "empty value: "+k)
			}
		}
	}
	completenessScore := 50
	if len(vars) > 0 && emptyCount > 0 {
		cut := (emptyCount * 50) / len(vars)
		if cut > 50 {
			cut = 50
		}
		completenessScore -= cut
		total -= cut
	}
	categories["completeness"] = completenessScore

	// Security checks
	securityPenalty := 0
	for k, v := range vars {
		for _, pat := range opts.SensitivePatterns {
			if strings.Contains(strings.ToUpper(k), pat) && strings.TrimSpace(v) == "" {
				securityPenalty++
				deductions = append(deductions, "sensitive key has empty value: "+k)
				break
			}
		}
	}
	securityScore := 25
	if securityPenalty > 0 {
		cut := securityPenalty * 10
		if cut > 25 {
			cut = 25
		}
		securityScore -= cut
		total -= cut
	}
	categories["security"] = securityScore

	if total < 0 {
		total = 0
	}
	return Result{Total: total, Categories: categories, Deductions: deductions}
}
