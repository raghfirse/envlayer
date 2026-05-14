// Package envfilter provides predicate-based filtering of environment variable maps.
// Filters can be composed to create complex selection criteria over env var keys and values.
package envfilter

import "strings"

// Predicate is a function that returns true if the given key-value pair should be included.
type Predicate func(key, value string) bool

// Filter returns a new map containing only entries for which all predicates return true.
// If no predicates are given, the original map is returned unchanged (shallow copy).
func Filter(vars map[string]string, predicates ...Predicate) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		if matchAll(k, v, predicates) {
			out[k] = v
		}
	}
	return out
}

// HasPrefix returns a Predicate that matches keys with the given prefix.
func HasPrefix(prefix string) Predicate {
	return func(key, _ string) bool {
		return strings.HasPrefix(key, prefix)
	}
}

// HasSuffix returns a Predicate that matches keys with the given suffix.
func HasSuffix(suffix string) Predicate {
	return func(key, _ string) bool {
		return strings.HasSuffix(key, suffix)
	}
}

// ValueNotEmpty returns a Predicate that excludes entries with empty values.
func ValueNotEmpty() Predicate {
	return func(_, value string) bool {
		return value != ""
	}
}

// KeyContains returns a Predicate that matches keys containing the given substring.
func KeyContains(sub string) Predicate {
	return func(key, _ string) bool {
		return strings.Contains(key, sub)
	}
}

// Not negates a predicate.
func Not(p Predicate) Predicate {
	return func(key, value string) bool {
		return !p(key, value)
	}
}

// Any returns a Predicate that is true when at least one of the given predicates is true.
func Any(predicates ...Predicate) Predicate {
	return func(key, value string) bool {
		for _, p := range predicates {
			if p(key, value) {
				return true
			}
		}
		return false
	}
}

func matchAll(key, value string, predicates []Predicate) bool {
	for _, p := range predicates {
		if !p(key, value) {
			return false
		}
	}
	return true
}
