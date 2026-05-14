// Package envaliaser provides key aliasing for environment variable maps.
// It allows renaming keys by defining alias mappings, supporting both
// one-to-one renames and fan-out (one key copied to multiple aliases).
package envaliaser

import "fmt"

// AliasMap maps original key names to one or more alias names.
type AliasMap map[string][]string

// Options controls alias behaviour.
type Options struct {
	// KeepOriginal retains the original key alongside the alias(es).
	KeepOriginal bool
	// FailOnMissing returns an error when a source key is absent.
	FailOnMissing bool
}

// DefaultOptions returns sensible defaults: drop original keys, ignore missing.
func DefaultOptions() Options {
	return Options{
		KeepOriginal:  false,
		FailOnMissing: false,
	}
}

// Apply rewrites vars according to aliases and opts.
// The returned map is always a new allocation; the input is never mutated.
func Apply(vars map[string]string, aliases AliasMap, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(vars))
	// Copy all entries first so non-aliased keys are preserved.
	for k, v := range vars {
		out[k] = v
	}

	for src, targets := range aliases {
		val, ok := vars[src]
		if !ok {
			if opts.FailOnMissing {
				return nil, fmt.Errorf("envaliaser: source key %q not found", src)
			}
			continue
		}
		for _, target := range targets {
			out[target] = val
		}
		if !opts.KeepOriginal {
			delete(out, src)
		}
	}
	return out, nil
}

// Invert builds a reverse mapping from each alias back to its source key.
// If an alias appears more than once the last source wins.
func Invert(aliases AliasMap) map[string]string {
	inv := make(map[string]string)
	for src, targets := range aliases {
		for _, t := range targets {
			inv[t] = src
		}
	}
	return inv
}
