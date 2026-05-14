// Package envaliaser provides key aliasing for environment variable maps.
//
// It allows callers to rename or fan-out environment variable keys by
// declaring an AliasMap — a mapping from source key names to one or more
// target (alias) names.
//
// Example usage:
//
//	aliases := envaliaser.AliasMap{
//		"DB_HOST": {"DATABASE_HOST", "PG_HOST"},
//	}
//	out, err := envaliaser.Apply(vars, aliases, envaliaser.DefaultOptions())
//
// By default the original key is removed after aliasing. Set
// Options.KeepOriginal = true to retain it. Set Options.FailOnMissing = true
// to return an error when a source key is absent from the input map.
//
// Invert builds the reverse mapping (alias -> source) which is useful for
// translating aliased output back to canonical key names.
package envaliaser
