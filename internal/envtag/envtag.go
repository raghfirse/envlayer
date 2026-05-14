// Package envtag provides utilities for tagging and filtering environment
// variables using structured key metadata embedded in key naming conventions
// or explicit tag maps.
package envtag

import (
	"fmt"
	"sort"
	"strings"
)

// Tag represents a label attached to one or more environment variable keys.
type Tag struct {
	Name string
	Keys []string
}

// Index maps each env key to a set of tag names.
type Index map[string][]string

// Build constructs an Index from a slice of Tags.
func Build(tags []Tag) Index {
	idx := make(Index)
	for _, t := range tags {
		for _, k := range t.Keys {
			idx[k] = append(idx[k], t.Name)
		}
	}
	return idx
}

// FilterByTag returns a new map containing only keys that carry the given tag.
func FilterByTag(vars map[string]string, idx Index, tag string) map[string]string {
	out := make(map[string]string)
	for k, v := range vars {
		for _, t := range idx[k] {
			if t == tag {
				out[k] = v
				break
			}
		}
	}
	return out
}

// TagsFor returns the sorted list of tags assigned to a given key.
func TagsFor(idx Index, key string) []string {
	tags := make([]string, len(idx[key]))
	copy(tags, idx[key])
	sort.Strings(tags)
	return tags
}

// AllTags returns a deduplicated, sorted list of all tag names in the index.
func AllTags(idx Index) []string {
	seen := make(map[string]struct{})
	for _, tags := range idx {
		for _, t := range tags {
			seen[t] = struct{}{}
		}
	}
	out := make([]string, 0, len(seen))
	for t := range seen {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

// Validate checks that every key in the index exists in vars, returning
// an error listing any unknown keys.
func Validate(vars map[string]string, idx Index) error {
	var missing []string
	for k := range idx {
		if _, ok := vars[k]; !ok {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("envtag: unknown keys in index: %s", strings.Join(missing, ", "))
	}
	return nil
}
