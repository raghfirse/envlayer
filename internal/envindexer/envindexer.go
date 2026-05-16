// Package envindexer builds a reverse index from env var values to their keys,
// enabling fast lookup of which keys hold a given value.
package envindexer

import "sort"

// Index maps each unique value to the sorted list of keys that hold it.
type Index map[string][]string

// Options controls how the index is built.
type Options struct {
	// CaseFoldValues normalises values to lowercase before indexing.
	CaseFoldValues bool
	// ExcludeEmpty skips keys whose value is the empty string.
	ExcludeEmpty bool
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		CaseFoldValues: false,
		ExcludeEmpty:   true,
	}
}

// Build constructs an Index from the provided env map.
func Build(vars map[string]string, opts Options) Index {
	idx := make(Index)

	for k, v := range vars {
		if opts.ExcludeEmpty && v == "" {
			continue
		}
		if opts.CaseFoldValues {
			v = toLower(v)
		}
		idx[v] = append(idx[v], k)
	}

	for v := range idx {
		sort.Strings(idx[v])
	}

	return idx
}

// Lookup returns the keys that hold the given value, or nil if none.
func Lookup(idx Index, value string) []string {
	keys, ok := idx[value]
	if !ok {
		return nil
	}
	out := make([]string, len(keys))
	copy(out, keys)
	return out
}

// Values returns all indexed values in sorted order.
func Values(idx Index) []string {
	out := make([]string, 0, len(idx))
	for v := range idx {
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

func toLower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}
