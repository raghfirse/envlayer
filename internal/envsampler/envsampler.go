// Package envsampler provides utilities for sampling a subset of environment
// variables from a map, supporting random sampling, top-N by key length, and
// deterministic sampling by key prefix pattern.
package envsampler

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

// Options controls how sampling is performed.
type Options struct {
	// Seed is used for reproducible random sampling. Zero means non-deterministic.
	Seed int64
	// Prefix filters keys before sampling when non-empty.
	Prefix string
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{}
}

// Sample returns up to n key-value pairs from vars, chosen randomly.
// If n >= len(vars), all entries are returned (in sorted key order).
func Sample(vars map[string]string, n int, opts Options) (map[string]string, error) {
	if n < 0 {
		return nil, fmt.Errorf("envsampler: n must be non-negative, got %d", n)
	}

	pool := filtered(vars, opts.Prefix)
	if n >= len(pool) {
		return toMap(pool), nil
	}

	r := rand.New(rand.NewSource(opts.Seed))
	r.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
	return toMap(pool[:n]), nil
}

// TopN returns the n entries whose keys are longest (ties broken alphabetically).
// If n >= len(vars), all entries are returned.
func TopN(vars map[string]string, n int, opts Options) (map[string]string, error) {
	if n < 0 {
		return nil, fmt.Errorf("envsampler: n must be non-negative, got %d", n)
	}

	pool := filtered(vars, opts.Prefix)
	sort.Slice(pool, func(i, j int) bool {
		if len(pool[i][0]) != len(pool[j][0]) {
			return len(pool[i][0]) > len(pool[j][0])
		}
		return pool[i][0] < pool[j][0]
	})

	if n >= len(pool) {
		return toMap(pool), nil
	}
	return toMap(pool[:n]), nil
}

// filtered returns key-value pairs from vars whose keys start with prefix.
// If prefix is empty, all pairs are included. Result is sorted by key.
func filtered(vars map[string]string, prefix string) [][2]string {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		if prefix == "" || strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	pairs := make([][2]string, len(keys))
	for i, k := range keys {
		pairs[i] = [2]string{k, vars[k]}
	}
	return pairs
}

func toMap(pairs [][2]string) map[string]string {
	out := make(map[string]string, len(pairs))
	for _, p := range pairs {
		out[p[0]] = p[1]
	}
	return out
}
