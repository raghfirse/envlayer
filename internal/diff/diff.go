// Package diff provides utilities for comparing two sets of environment
// variables and reporting additions, removals, and changed values.
package diff

import "sort"

// ChangeType describes the kind of change detected for a key.
type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Changed ChangeType = "changed"
)

// Entry represents a single difference between two env maps.
type Entry struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Compare returns the ordered list of differences between the base and next
// environment maps. Keys present only in base are Removed; keys present only
// in next are Added; keys present in both with differing values are Changed.
func Compare(base, next map[string]string) []Entry {
	seen := make(map[string]bool)
	var entries []Entry

	for k, oldVal := range base {
		seen[k] = true
		if newVal, ok := next[k]; !ok {
			entries = append(entries, Entry{
				Key:      k,
				Type:     Removed,
				OldValue: oldVal,
			})
		} else if newVal != oldVal {
			entries = append(entries, Entry{
				Key:      k,
				Type:     Changed,
				OldValue: oldVal,
				NewValue: newVal,
			})
		}
	}

	for k, newVal := range next {
		if !seen[k] {
			entries = append(entries, Entry{
				Key:      k,
				Type:     Added,
				NewValue: newVal,
			})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return entries
}
