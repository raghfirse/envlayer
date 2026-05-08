// Package differ provides utilities for computing and formatting
// environment variable diffs between two named snapshots or maps.
package differ

import (
	"fmt"
	"sort"
)

// ChangeKind describes the type of change for a single key.
type ChangeKind string

const (
	Added   ChangeKind = "added"
	Removed ChangeKind = "removed"
	Changed ChangeKind = "changed"
)

// Change represents a single variable change between two environments.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// Result holds the full diff between two env maps.
type Result struct {
	From    string
	To      string
	Changes []Change
}

// Diff computes the difference between two env maps, labelled from/to.
func Diff(from, to map[string]string, fromLabel, toLabel string) Result {
	result := Result{From: fromLabel, To: toLabel}

	keys := unionKeys(from, to)
	for _, k := range keys {
		oldVal, inFrom := from[k]
		newVal, inTo := to[k]

		switch {
		case inFrom && !inTo:
			result.Changes = append(result.Changes, Change{Key: k, Kind: Removed, OldValue: oldVal})
		case !inFrom && inTo:
			result.Changes = append(result.Changes, Change{Key: k, Kind: Added, NewValue: newVal})
		case oldVal != newVal:
			result.Changes = append(result.Changes, Change{Key: k, Kind: Changed, OldValue: oldVal, NewValue: newVal})
		}
	}
	return result
}

// Summary returns a human-readable one-line summary of the diff result.
func Summary(r Result) string {
	added, removed, changed := 0, 0, 0
	for _, c := range r.Changes {
		switch c.Kind {
		case Added:
			added++
		case Removed:
			removed++
		case Changed:
			changed++
		}
	}
	return fmt.Sprintf("%s → %s: +%d -%d ~%d", r.From, r.To, added, removed, changed)
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
