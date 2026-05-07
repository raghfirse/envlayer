// Package transformer provides key/value transformation utilities
// for environment variable maps, such as prefix injection, key casing,
// and value trimming.
package transformer

import (
	"strings"
)

// Options controls which transformations are applied.
type Options struct {
	// AddPrefix prepends a string to every key.
	AddPrefix string

	// StripPrefix removes a leading string from every key.
	StripPrefix string

	// UppercaseKeys converts all keys to UPPER_CASE.
	UppercaseKeys bool

	// LowercaseKeys converts all keys to lower_case.
	LowercaseKeys bool

	// TrimValues strips leading/trailing whitespace from values.
	TrimValues bool
}

// Transform applies the given Options to a copy of vars and returns the
// transformed map. The original map is never mutated.
func Transform(vars map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		if opts.StripPrefix != "" {
			k = strings.TrimPrefix(k, opts.StripPrefix)
		}
		if opts.AddPrefix != "" {
			k = opts.AddPrefix + k
		}
		if opts.UppercaseKeys {
			k = strings.ToUpper(k)
		} else if opts.LowercaseKeys {
			k = strings.ToLower(k)
		}
		if opts.TrimValues {
			v = strings.TrimSpace(v)
		}
		out[k] = v
	}
	return out
}
