// Package envfreezer provides functionality to freeze (lock) a set of
// environment variables, preventing any keys from being modified or removed
// in subsequent merge operations. Frozen keys can only be overridden
// explicitly by passing AllowOverride.
package envfreezer

import "fmt"

// Options configures the behaviour of the Freeze operation.
type Options struct {
	// AllowOverride lists keys that are exempt from the freeze and may be
	// changed freely even after freezing.
	AllowOverride []string

	// StrictMode causes Freeze to return an error when a frozen key would
	// be mutated, rather than silently keeping the original value.
	StrictMode bool
}

// DefaultOptions returns an Options value with sensible defaults.
func DefaultOptions() Options {
	return Options{}
}

// Freeze takes a frozen snapshot of base and applies updates on top of it.
// Keys present in base are protected: if updates tries to change or delete
// them, the original value is kept (or an error is returned in StrictMode).
// New keys that exist only in updates are always added.
func Freeze(base, updates map[string]string, opts Options) (map[string]string, error) {
	exempt := make(map[string]bool, len(opts.AllowOverride))
	for _, k := range opts.AllowOverride {
		exempt[k] = true
	}

	result := make(map[string]string, len(base)+len(updates))

	// Copy base into result — these values are frozen.
	for k, v := range base {
		result[k] = v
	}

	// Apply updates, respecting the freeze.
	for k, v := range updates {
		original, frozen := base[k]
		if frozen && !exempt[k] {
			if v != original && opts.StrictMode {
				return nil, fmt.Errorf("envfreezer: key %q is frozen and cannot be changed", k)
			}
			// Keep original value; do not apply update.
			continue
		}
		result[k] = v
	}

	return result, nil
}

// FrozenKeys returns the sorted list of keys that are currently frozen
// (i.e. present in base and not in the AllowOverride list).
func FrozenKeys(base map[string]string, opts Options) []string {
	exempt := make(map[string]bool, len(opts.AllowOverride))
	for _, k := range opts.AllowOverride {
		exempt[k] = true
	}

	keys := make([]string, 0, len(base))
	for k := range base {
		if !exempt[k] {
			keys = append(keys, k)
		}
	}
	sortStrings(keys)
	return keys
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
