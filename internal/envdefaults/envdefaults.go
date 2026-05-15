// Package envdefaults provides utilities for applying default values to
// environment variable maps. Keys that are missing or have empty values
// can be populated from a defaults map, with configurable behaviour.
package envdefaults

import "fmt"

// Options controls how defaults are applied.
type Options struct {
	// OverwriteEmpty replaces existing keys whose value is the empty string.
	OverwriteEmpty bool
	// FailOnConflict returns an error when a non-empty key already exists in
	// the target and a default would otherwise be skipped.
	FailOnConflict bool
}

// DefaultOptions returns a sensible out-of-the-box configuration.
func DefaultOptions() Options {
	return Options{
		OverwriteEmpty: true,
		FailOnConflict: false,
	}
}

// Applied records which keys were set and which were skipped.
type Applied struct {
	Set     []string
	Skipped []string
}

// Apply merges defaults into target according to opts.
// target is never mutated; a new map is returned alongside a report.
func Apply(target, defaults map[string]string, opts Options) (map[string]string, Applied, error) {
	out := make(map[string]string, len(target))
	for k, v := range target {
		out[k] = v
	}

	var report Applied

	for k, defVal := range defaults {
		existing, exists := out[k]
		switch {
		case !exists || (existing == "" && opts.OverwriteEmpty):
			out[k] = defVal
			report.Set = append(report.Set, k)
		case opts.FailOnConflict && existing != "" && existing != defVal:
			return nil, Applied{}, fmt.Errorf("envdefaults: conflict on key %q: existing %q vs default %q", k, existing, defVal)
		default:
			report.Skipped = append(report.Skipped, k)
		}
	}

	sortStrings(report.Set)
	sortStrings(report.Skipped)
	return out, report, nil
}

// ApplySimple is a convenience wrapper that uses DefaultOptions and panics on
// error (useful in tests and one-shot CLI paths).
func ApplySimple(target, defaults map[string]string) map[string]string {
	out, _, err := Apply(target, defaults, DefaultOptions())
	if err != nil {
		panic(err)
	}
	return out
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
