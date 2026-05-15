// Package envtrimmer provides utilities for trimming environment variable
// maps by removing, capping, or filtering entries based on configurable rules.
package envtrimmer

import "strings"

// Options controls the behaviour of Trim.
type Options struct {
	// MaxValueLen caps each value to this many characters (0 = unlimited).
	MaxValueLen int

	// OmitEmpty removes keys whose value is the empty string.
	OmitEmpty bool

	// OmitPrefixes removes any key that starts with one of these prefixes.
	OmitPrefixes []string

	// OmitKeys removes these exact keys.
	OmitKeys []string
}

// DefaultOptions returns an Options with sensible defaults (no trimming).
func DefaultOptions() Options {
	return Options{}
}

// Trim returns a new map derived from vars with the supplied options applied.
// The original map is never mutated.
func Trim(vars map[string]string, opts Options) map[string]string {
	omitKeySet := make(map[string]struct{}, len(opts.OmitKeys))
	for _, k := range opts.OmitKeys {
		omitKeySet[k] = struct{}{}
	}

	out := make(map[string]string, len(vars))
	for k, v := range vars {
		// Drop explicitly omitted keys.
		if _, skip := omitKeySet[k]; skip {
			continue
		}

		// Drop keys matching any omit-prefix.
		if hasAnyPrefix(k, opts.OmitPrefixes) {
			continue
		}

		// Drop empty values when requested.
		if opts.OmitEmpty && v == "" {
			continue
		}

		// Cap value length.
		if opts.MaxValueLen > 0 && len(v) > opts.MaxValueLen {
			v = v[:opts.MaxValueLen]
		}

		out[k] = v
	}
	return out
}

// TrimKeys returns a copy of vars containing only the supplied keys.
func TrimKeys(vars map[string]string, keep []string) map[string]string {
	out := make(map[string]string, len(keep))
	for _, k := range keep {
		if v, ok := vars[k]; ok {
			out[k] = v
		}
	}
	return out
}

func hasAnyPrefix(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}
