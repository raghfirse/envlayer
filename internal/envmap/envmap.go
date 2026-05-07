package envmap

import "sort"

// Entry represents a single environment variable with its key and value.
type Entry struct {
	Key   string
	Value string
}

// FromMap converts a map[string]string into a sorted slice of Entry.
func FromMap(m map[string]string) []Entry {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]Entry, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, Entry{Key: k, Value: m[k]})
	}
	return entries
}

// ToMap converts a slice of Entry into a map[string]string.
// Later entries overwrite earlier ones if keys collide.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

// Filter returns only entries whose keys satisfy the predicate.
func Filter(entries []Entry, predicate func(key string) bool) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if predicate(e.Key) {
			out = append(out, e)
		}
	}
	return out
}

// MapValues returns a new slice with each value transformed by fn.
func MapValues(entries []Entry, fn func(key, value string) string) []Entry {
	out := make([]Entry, len(entries))
	for i, e := range entries {
		out[i] = Entry{Key: e.Key, Value: fn(e.Key, e.Value)}
	}
	return out
}
