// Package envpatcher provides a lightweight patch mechanism for environment
// variable maps.
//
// A patch is an ordered list of [Op] values, each describing one of three
// mutation kinds:
//
//   - OpSet    – create or overwrite a key with a given value.
//   - OpDelete – remove a key; silently skipped when the key is absent.
//   - OpRename – move a key to a new name, preserving its value.
//
// [Apply] executes all ops against a copy of the source map and returns a
// [Result] that contains the patched map together with per-op audit logs
// (Applied / Skipped). The original map is never mutated.
//
// Example:
//
//	ops := []envpatcher.Op{
//		{Kind: envpatcher.OpSet,    Key: "LOG_LEVEL", Value: "info"},
//		{Kind: envpatcher.OpDelete, Key: "LEGACY_FLAG"},
//		{Kind: envpatcher.OpRename, Key: "DB_PASS", NewKey: "DATABASE_PASSWORD"},
//	}
//	res, err := envpatcher.Apply(vars, ops)
package envpatcher
