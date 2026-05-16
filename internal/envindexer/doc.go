// Package envindexer provides a reverse-index over an environment variable map.
//
// It answers the question: "which keys share this value?" — useful for
// detecting duplicate values, auditing sensitive data propagation, and
// building cross-reference reports across layered .env files.
//
// # Building an index
//
//	idx := envindexer.Build(vars, envindexer.DefaultOptions())
//
// # Looking up keys by value
//
//	keys := envindexer.Lookup(idx, "localhost")
//
// # Listing all indexed values
//
//	vals := envindexer.Values(idx)
package envindexer
