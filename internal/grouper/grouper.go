// Package grouper organises a flat map of environment variables into
// named groups based on a key-prefix convention.
//
// For example, given the prefix separator "_", the key "DB_HOST" belongs
// to the group "DB" and the key "APP_PORT" belongs to the group "APP".
// Keys that contain no separator are placed into the special "" (root) group.
package grouper

import (
	"sort"
	"strings"
)

// Group holds the variables that share a common prefix.
type Group struct {
	// Name is the prefix that all keys in this group share (may be empty
	// for root-level keys).
	Name string

	// Vars contains the key/value pairs belonging to this group.  The keys
	// are stored without the leading prefix and separator.
	Vars map[string]string
}

// GroupBy splits vars into groups using sep as the separator between the
// prefix and the rest of the key.  Only the first occurrence of sep is
// used, so "DB_HOST_PORT" is placed into group "DB" with key "HOST_PORT".
//
// When sep is empty, GroupBy returns a single root group containing all
// variables unchanged.
func GroupBy(vars map[string]string, sep string) []Group {
	if sep == "" {
		return []Group{{Name: "", Vars: copyMap(vars)}}
	}

	buckets := map[string]map[string]string{}

	for k, v := range vars {
		prefix, rest, found := strings.Cut(k, sep)
		if !found {
			// No separator – root group.
			prefix = ""
			rest = k
		}
		if buckets[prefix] == nil {
			buckets[prefix] = map[string]string{}
		}
		buckets[prefix][rest] = v
	}

	groups := make([]Group, 0, len(buckets))
	for name, m := range buckets {
		groups = append(groups, Group{Name: name, Vars: m})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	return groups
}

// Flatten is the inverse of GroupBy: it reconstructs a flat map from a
// slice of Groups, re-attaching the group name as a prefix separated by sep.
//
// Root-group keys (Group.Name == "") are written without any prefix.
func Flatten(groups []Group, sep string) map[string]string {
	out := map[string]string{}
	for _, g := range groups {
		for k, v := range g.Vars {
			key := k
			if g.Name != "" {
				key = g.Name + sep + k
			}
			out[key] = v
		}
	}
	return out
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
