// Package envexpander expands shorthand or abbreviated environment variable
// keys to their full canonical names using a user-supplied expansion map.
package envexpander

import "fmt"

// Options controls the behaviour of the expander.
type Options struct {
	// FailOnUnknown causes Expand to return an error when a key has no
	// expansion entry and is not already a full-length key.
	FailOnUnknown bool

	// KeepOriginal retains the original short key alongside the expanded one.
	KeepOriginal bool
}

// DefaultOptions returns a safe default configuration.
func DefaultOptions() Options {
	return Options{
		FailOnUnknown: false,
		KeepOriginal:  false,
	}
}

// Expand iterates over vars and replaces any key found in expansions with its
// canonical name. Keys not present in expansions are kept as-is unless
// FailOnUnknown is set. If KeepOriginal is true the original key is also
// retained in the output map.
func Expand(vars map[string]string, expansions map[string]string, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(vars))

	for k, v := range vars {
		canonical, ok := expansions[k]
		if !ok {
			if opts.FailOnUnknown {
				return nil, fmt.Errorf("envexpander: no expansion defined for key %q", k)
			}
			out[k] = v
			continue
		}

		out[canonical] = v
		if opts.KeepOriginal && k != canonical {
			out[k] = v
		}
	}

	return out, nil
}

// Invert reverses an expansion map so that canonical names map back to their
// shorthand equivalents. If multiple shorthands share the same canonical name
// the last one wins.
func Invert(expansions map[string]string) map[string]string {
	inv := make(map[string]string, len(expansions))
	for short, canonical := range expansions {
		inv[canonical] = short
	}
	return inv
}
