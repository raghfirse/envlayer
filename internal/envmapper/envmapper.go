// Package envmapper provides utilities for remapping environment variable
// keys and values using a declarative mapping definition.
package envmapper

import (
	"fmt"
	"sort"
)

// Rule defines a single mapping operation.
type Rule struct {
	From string // source key
	To   string // destination key
	Keep bool   // if true, retain the original key as well
}

// Options controls the behaviour of Apply.
type Options struct {
	Rules          []Rule
	FailOnMissing  bool // return error when a source key is absent
	DropUnmapped   bool // drop all keys that are not referenced by any rule
}

// DefaultOptions returns a safe default configuration.
func DefaultOptions() Options {
	return Options{
		FailOnMissing: false,
		DropUnmapped:  false,
	}
}

// Apply remaps vars according to opts.Rules and returns a new map.
func Apply(vars map[string]string, opts Options) (map[string]string, error) {
	mapped := make(map[string]string, len(vars))

	// Track which source keys were consumed by a rule.
	consumed := make(map[string]bool)

	for _, rule := range opts.Rules {
		if rule.From == "" || rule.To == "" {
			return nil, fmt.Errorf("envmapper: rule has empty From or To field")
		}
		val, ok := vars[rule.From]
		if !ok {
			if opts.FailOnMissing {
				return nil, fmt.Errorf("envmapper: source key %q not found", rule.From)
			}
			continue
		}
		mapped[rule.To] = val
		consumed[rule.From] = true
		if rule.Keep {
			mapped[rule.From] = val
		}
	}

	if !opts.DropUnmapped {
		for k, v := range vars {
			if !consumed[k] {
				mapped[k] = v
			}
		}
	}

	return mapped, nil
}

// Keys returns a sorted slice of all destination keys produced by rules.
func Keys(opts Options) []string {
	seen := make(map[string]bool)
	for _, r := range opts.Rules {
		if r.To != "" {
			seen[r.To] = true
		}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
