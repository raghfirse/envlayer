// Package envstats provides statistical analysis over a set of environment variables,
// including counts, key length distribution, and value pattern summaries.
package envstats

import (
	"sort"
	"strings"
)

// Stats holds aggregate information about an env map.
type Stats struct {
	TotalKeys      int
	EmptyValues    int
	SensitiveKeys  int
	AvgKeyLength   float64
	AvgValueLength float64
	PrefixGroups   map[string]int // top-level prefix (before first "_") -> count
}

// sensitivePatterns are substrings that indicate a key may be sensitive.
var sensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "KEY", "PRIVATE", "CREDENTIAL",
}

// Compute derives statistics from the provided env map.
func Compute(vars map[string]string) Stats {
	if len(vars) == 0 {
		return Stats{PrefixGroups: map[string]int{}}
	}

	var totalKeyLen, totalValLen int
	prefixGroups := map[string]int{}
	empty := 0
	sensitive := 0

	for k, v := range vars {
		totalKeyLen += len(k)
		totalValLen += len(v)

		if v == "" {
			empty++
		}

		if isSensitive(k) {
			sensitive++
		}

		prefix := topPrefix(k)
		prefixGroups[prefix]++
	}

	n := len(vars)
	return Stats{
		TotalKeys:      n,
		EmptyValues:    empty,
		SensitiveKeys:  sensitive,
		AvgKeyLength:   float64(totalKeyLen) / float64(n),
		AvgValueLength: float64(totalValLen) / float64(n),
		PrefixGroups:   prefixGroups,
	}
}

// TopPrefixes returns prefix group names sorted by descending count.
func TopPrefixes(s Stats) []string {
	type kv struct {
		key   string
		count int
	}
	pairs := make([]kv, 0, len(s.PrefixGroups))
	for k, c := range s.PrefixGroups {
		pairs = append(pairs, kv{k, c})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].key < pairs[j].key
	})
	out := make([]string, len(pairs))
	for i, p := range pairs {
		out[i] = p.key
	}
	return out
}

func topPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}

func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}
