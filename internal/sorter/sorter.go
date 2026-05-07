// Package sorter provides utilities for ordering environment variable maps
// by key name, value length, or insertion-defined priority lists.
package sorter

import (
	"sort"
)

// Order defines the sort direction.
type Order int

const (
	Ascending  Order = iota
	Descending
)

// ByKey returns a new map with keys sorted alphabetically into a slice of
// key-value pairs, preserving all entries from vars.
func ByKey(vars map[string]string, order Order) []Entry {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if order == Descending {
		for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
			keys[i], keys[j] = keys[j], keys[i]
		}
	}
	entries := make([]Entry, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, Entry{Key: k, Value: vars[k]})
	}
	return entries
}

// ByPriority returns entries ordered so that keys listed in priority appear
// first (in priority order), followed by remaining keys sorted alphabetically.
func ByPriority(vars map[string]string, priority []string) []Entry {
	seen := make(map[string]bool, len(priority))
	entries := make([]Entry, 0, len(vars))

	for _, k := range priority {
		if v, ok := vars[k]; ok {
			entries = append(entries, Entry{Key: k, Value: v})
			seen[k] = true
		}
	}

	remaining := make([]string, 0)
	for k := range vars {
		if !seen[k] {
			remaining = append(remaining, k)
		}
	}
	sort.Strings(remaining)
	for _, k := range remaining {
		entries = append(entries, Entry{Key: k, Value: vars[k]})
	}
	return entries
}

// ToMap converts a slice of Entry back into a plain map.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

// Entry holds a single key-value pair.
type Entry struct {
	Key   string
	Value string
}
