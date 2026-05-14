// Package envmerger provides utilities for merging multiple environment
// variable maps with configurable conflict resolution strategies.
package envmerger

import "sort"

// Strategy defines how conflicts are resolved when merging env maps.
type Strategy int

const (
	// LastWins uses the value from the last map that defines a key.
	LastWins Strategy = iota
	// FirstWins uses the value from the first map that defines a key.
	FirstWins
	// ErrorOnConflict returns an error if the same key appears in multiple maps.
	ErrorOnConflict
)

// Options configures the merge behaviour.
type Options struct {
	Strategy Strategy
}

// DefaultOptions returns sensible merge defaults (LastWins).
func DefaultOptions() Options {
	return Options{Strategy: LastWins}
}

// Result holds the merged map and metadata about the merge.
type Result struct {
	Vars      map[string]string
	Conflicts []Conflict
}

// Conflict records a key that appeared in more than one source layer.
type Conflict struct {
	Key    string
	Values []string // values in layer order
}

// Merge combines the provided layers according to opts.
// Layers are applied in order; index 0 is the lowest-priority layer.
func Merge(layers []map[string]string, opts Options) (Result, error) {
	out := make(map[string]string)
	conflictIndex := make(map[string][]string)

	for _, layer := range layers {
		for _, k := range sortedKeys(layer) {
			v := layer[k]
			if existing, ok := out[k]; ok {
				if opts.Strategy == ErrorOnConflict && existing != v {
					return Result{}, &ConflictError{Key: k}
				}
				conflictIndex[k] = append(conflictIndex[k], v)
				if opts.Strategy == LastWins {
					out[k] = v
				}
			} else {
				out[k] = v
				conflictIndex[k] = []string{v}
			}
		}
	}

	var conflicts []Conflict
	for _, k := range sortedKeys(conflictIndex) {
		if len(conflictIndex[k]) > 1 {
			conflicts = append(conflicts, Conflict{Key: k, Values: conflictIndex[k]})
		}
	}

	return Result{Vars: out, Conflicts: conflicts}, nil
}

// ConflictError is returned when ErrorOnConflict strategy detects a duplicate key.
type ConflictError struct {
	Key string
}

func (e *ConflictError) Error() string {
	return "envmerger: conflict on key " + e.Key
}

func sortedKeys[M ~map[string]V, V any](m M) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
