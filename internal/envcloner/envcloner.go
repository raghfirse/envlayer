// Package envcloner provides utilities for deep-copying and cloning
// environment variable maps with optional transformations applied during cloning.
package envcloner

import "strings"

// Options controls how a clone operation is performed.
type Options struct {
	// KeyPrefix, if non-empty, is prepended to every key in the cloned map.
	KeyPrefix string
	// StripPrefix, if non-empty, is stripped from each key before cloning.
	// Stripping happens before KeyPrefix is applied.
	StripPrefix string
	// UppercaseKeys converts all keys to uppercase in the clone.
	UppercaseKeys bool
	// OmitEmpty excludes keys whose values are empty strings.
	OmitEmpty bool
}

// DefaultOptions returns an Options with no transformations.
func DefaultOptions() Options {
	return Options{}
}

// Clone returns a deep copy of src, applying any transformations specified in
// opts. The original map is never mutated.
func Clone(src map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		if opts.OmitEmpty && v == "" {
			continue
		}
		key := k
		if opts.StripPrefix != "" {
			key = strings.TrimPrefix(key, opts.StripPrefix)
		}
		if opts.UppercaseKeys {
			key = strings.ToUpper(key)
		}
		if opts.KeyPrefix != "" {
			key = opts.KeyPrefix + key
		}
		out[key] = v
	}
	return out
}

// CloneKeys returns a new map containing only the specified keys from src.
// Missing keys are silently skipped.
func CloneKeys(src map[string]string, keys []string) map[string]string {
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		if v, ok := src[k]; ok {
			out[k] = v
		}
	}
	return out
}

// CloneExclude returns a deep copy of src with the given keys omitted.
func CloneExclude(src map[string]string, exclude []string) map[string]string {
	skip := make(map[string]struct{}, len(exclude))
	for _, k := range exclude {
		skip[k] = struct{}{}
	}
	out := make(map[string]string, len(src))
	for k, v := range src {
		if _, found := skip[k]; !found {
			out[k] = v
		}
	}
	return out
}
