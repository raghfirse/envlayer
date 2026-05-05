// Package validator provides utilities for validating environment variable
// maps, checking for required keys and detecting potential issues.
package validator

import (
	"fmt"
	"strings"
)

// Result holds the outcome of a validation pass.
type Result struct {
	Missing  []string
	Warnings []string
}

// IsValid returns true when there are no missing required keys.
func (r Result) IsValid() bool {
	return len(r.Missing) == 0
}

// Error returns a formatted error string listing missing keys, or empty string
// when the result is valid.
func (r Result) Error() string {
	if r.IsValid() {
		return ""
	}
	return fmt.Sprintf("missing required keys: %s", strings.Join(r.Missing, ", "))
}

// Validate checks that every key listed in required is present and non-empty
// in env. It also appends a warning for any key whose value is an empty string
// even when the key is not required.
func Validate(env map[string]string, required []string) Result {
	res := Result{}

	requiredSet := make(map[string]struct{}, len(required))
	for _, k := range required {
		requiredSet[k] = struct{}{}
	}

	for _, k := range required {
		v, ok := env[k]
		if !ok || strings.TrimSpace(v) == "" {
			res.Missing = append(res.Missing, k)
		}
	}

	for k, v := range env {
		if _, isRequired := requiredSet[k]; isRequired {
			continue
		}
		if strings.TrimSpace(v) == "" {
			res.Warnings = append(res.Warnings, fmt.Sprintf("key %q is present but empty", k))
		}
	}

	return res
}
