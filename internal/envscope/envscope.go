// Package envscope provides scoped views over environment variable maps,
// allowing consumers to work with a filtered, namespaced subset of variables.
package envscope

import (
	"fmt"
	"sort"
	"strings"
)

// Scope represents a named, filtered view over a set of environment variables.
type Scope struct {
	Name   string
	vars   map[string]string
	prefix string
}

// New creates a Scope from the given vars, filtering keys that start with prefix.
// The prefix is stripped from keys within the scope.
func New(name, prefix string, vars map[string]string) *Scope {
	scoped := make(map[string]string)
	for k, v := range vars {
		if strings.HasPrefix(k, prefix) {
			stripped := strings.TrimPrefix(k, prefix)
			if stripped != "" {
				scoped[stripped] = v
			}
		}
	}
	return &Scope{Name: name, vars: scoped, prefix: prefix}
}

// Get returns the value for key within the scope, and whether it was found.
func (s *Scope) Get(key string) (string, bool) {
	v, ok := s.vars[key]
	return v, ok
}

// All returns a copy of all variables within the scope.
func (s *Scope) All() map[string]string {
	out := make(map[string]string, len(s.vars))
	for k, v := range s.vars {
		out[k] = v
	}
	return out
}

// Keys returns the sorted list of keys in the scope.
func (s *Scope) Keys() []string {
	keys := make([]string, 0, len(s.vars))
	for k := range s.vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Qualify returns the fully-qualified (prefixed) key for use outside the scope.
func (s *Scope) Qualify(key string) string {
	return s.prefix + key
}

// String returns a human-readable summary of the scope.
func (s *Scope) String() string {
	return fmt.Sprintf("Scope(%q, prefix=%q, keys=%d)", s.Name, s.prefix, len(s.vars))
}
