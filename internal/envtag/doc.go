// Package envtag provides tagging and filtering capabilities for environment
// variable maps.
//
// Tags are named labels that can be associated with one or more environment
// variable keys. Once an Index is built from a set of Tag definitions, the
// index can be used to:
//
//   - Filter a variable map to only keys carrying a specific tag.
//   - Enumerate all tags assigned to a given key.
//   - List all distinct tag names present in the index.
//   - Validate that every indexed key exists in a given variable map.
//
// Example usage:
//
//	tags := []envtag.Tag{
//		{Name: "secret", Keys: []string{"DB_PASSWORD", "API_KEY"}},
//	}
//	idx := envtag.Build(tags)
//	secrets := envtag.FilterByTag(vars, idx, "secret")
package envtag
