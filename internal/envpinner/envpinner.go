// Package envpinner provides functionality for pinning environment variable
// values so they cannot be changed by subsequent merge or override operations.
package envpinner

import "fmt"

// Options controls the behaviour of the Pin operation.
type Options struct {
	// StrictMode causes Pin to return an error when an attempt is made to
	// override a pinned key with a different value.
	StrictMode bool
	// Silent suppresses the collection of pin violation messages when
	// StrictMode is false.
	Silent bool
}

// DefaultOptions returns a safe default Options value.
func DefaultOptions() Options {
	return Options{StrictMode: false, Silent: false}
}

// Result is the outcome of a Pin operation.
type Result struct {
	// Vars is the final merged map with pinned values enforced.
	Vars map[string]string
	// Violations lists keys whose incoming value was discarded because the key
	// was pinned.
	Violations []string
}

// Pin merges incoming into base while keeping the values for any key listed in
// pinned unchanged. Keys in incoming that are not pinned are merged normally
// (incoming wins). Keys not present in either map are ignored.
func Pin(base, incoming map[string]string, pinned []string, opts Options) (Result, error) {
	pinnedSet := make(map[string]struct{}, len(pinned))
	for _, k := range pinned {
		pinnedSet[k] = struct{}{}
	}

	out := make(map[string]string, len(base))
	for k, v := range base {
		out[k] = v
	}

	var violations []string

	for k, v := range incoming {
		if _, isPinned := pinnedSet[k]; isPinned {
			existing, exists := out[k]
			if exists && existing != v {
				if opts.StrictMode {
					return Result{}, fmt.Errorf("envpinner: key %q is pinned and cannot be overridden", k)
				}
				if !opts.Silent {
					violations = append(violations, k)
				}
			}
			// Keep the pinned (base) value — do not overwrite.
			continue
		}
		out[k] = v
	}

	return Result{Vars: out, Violations: violations}, nil
}

// PinnedKeys returns the subset of keys from vars that appear in the pinned
// list, preserving the order of the pinned list.
func PinnedKeys(vars map[string]string, pinned []string) []string {
	var found []string
	for _, k := range pinned {
		if _, ok := vars[k]; ok {
			found = append(found, k)
		}
	}
	return found
}
