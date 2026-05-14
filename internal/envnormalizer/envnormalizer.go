// Package envnormalizer provides utilities for normalizing environment variable
// maps by applying consistent formatting rules such as trimming whitespace,
// collapsing empty values, and standardizing key casing.
package envnormalizer

import (
	"strings"
)

// Options controls which normalization steps are applied.
type Options struct {
	// TrimValues removes leading and trailing whitespace from all values.
	TrimValues bool
	// UppercaseKeys converts all keys to uppercase.
	UppercaseKeys bool
	// LowercaseKeys converts all keys to lowercase. Ignored if UppercaseKeys is true.
	LowercaseKeys bool
	// RemoveEmpty drops entries whose value is empty after trimming.
	RemoveEmpty bool
	// CollapseWhitespace replaces runs of whitespace within values with a single space.
	CollapseWhitespace bool
}

// DefaultOptions returns an Options with the most common safe defaults.
func DefaultOptions() Options {
	return Options{
		TrimValues: true,
		UppercaseKeys: false,
		LowercaseKeys: false,
		RemoveEmpty: false,
		CollapseWhitespace: false,
	}
}

// Normalize applies the given Options to src and returns a new normalized map.
// The original map is never mutated.
func Normalize(src map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		if opts.TrimValues {
			v = strings.TrimSpace(v)
		}
		if opts.CollapseWhitespace {
			v = collapseWS(v)
		}
		if opts.RemoveEmpty && v == "" {
			continue
		}
		if opts.UppercaseKeys {
			k = strings.ToUpper(k)
		} else if opts.LowercaseKeys {
			k = strings.ToLower(k)
		}
		out[k] = v
	}
	return out
}

// collapseWS replaces any run of whitespace characters with a single space.
func collapseWS(s string) string {
	var b strings.Builder
	inSpace := false
	for _, r := range s {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			if !inSpace {
				b.WriteRune(' ')
				inSpace = true
			}
		} else {
			b.WriteRune(r)
			inSpace = false
		}
	}
	return b.String()
}
