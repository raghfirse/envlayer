// Package envtrimmer provides utilities for trimming and pruning environment
// variable maps according to configurable rules.
//
// # Overview
//
// Trim applies a set of reduction rules to a map[string]string and returns a
// new, independent copy — the original is never modified.
//
// Supported rules (via Options):
//
//   - OmitEmpty      — drop keys whose value is the empty string.
//   - OmitKeys       — drop specific named keys.
//   - OmitPrefixes   — drop any key that starts with one of the given prefixes.
//   - MaxValueLen    — cap every value to at most N characters.
//
// TrimKeys is a convenience function that returns a copy containing only the
// explicitly listed keys, silently ignoring any that are absent.
//
// # Example
//
//	opts := envtrimmer.DefaultOptions()
//	opts.OmitEmpty = true
//	opts.OmitPrefixes = []string{"INTERNAL_"}
//	clean := envtrimmer.Trim(vars, opts)
package envtrimmer
