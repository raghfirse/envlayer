// Package envrewriter rewrites environment variable maps by applying
// find-and-replace rules to keys, values, or both.
package envrewriter

import "strings"

// Rule describes a single rewrite operation.
type Rule struct {
	// Target selects what to rewrite: "key", "value", or "both".
	Target string
	// Find is the substring to search for.
	Find string
	// Replace is the replacement string.
	Replace string
}

// Options controls Rewrite behaviour.
type Options struct {
	// Rules is the ordered list of rewrite rules to apply.
	Rules []Rule
	// CaseSensitive controls whether Find matching is case-sensitive.
	CaseSensitive bool
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{CaseSensitive: true}
}

// Rewrite applies all rules in opts to vars and returns a new map.
// Original keys that are renamed retain their values unless a value rule
// also matches. If two rules rename the same key the last writer wins.
func Rewrite(vars map[string]string, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(vars))

	for k, v := range vars {
		newKey := k
		newVal := v

		for _, r := range opts.Rules {
			find := r.Find
			if find == "" {
				continue
			}

			applyKey := r.Target == "key" || r.Target == "both"
			applyVal := r.Target == "value" || r.Target == "both"

			if applyKey {
				newKey = replaceStr(newKey, find, r.Replace, opts.CaseSensitive)
			}
			if applyVal {
				newVal = replaceStr(newVal, find, r.Replace, opts.CaseSensitive)
			}
		}

		out[newKey] = newVal
	}

	return out, nil
}

func replaceStr(s, find, replace string, caseSensitive bool) string {
	if caseSensitive {
		return strings.ReplaceAll(s, find, replace)
	}
	// Case-insensitive replacement via manual scan.
	lower := strings.ToLower(s)
	lowerFind := strings.ToLower(find)
	var result strings.Builder
	for i := 0; i < len(s); {
		idx := strings.Index(lower[i:], lowerFind)
		if idx == -1 {
			result.WriteString(s[i:])
			break
		}
		result.WriteString(s[i : i+idx])
		result.WriteString(replace)
		i += idx + len(find)
	}
	return result.String()
}
