// Package envrotator applies controlled rotations to environment variable maps.
//
// A rotation replaces the value of one or more named keys with new values.
// The original map is never modified; Rotate always returns a fresh copy.
//
// Each call to Rotate produces a Result that contains:
//   - Vars: the updated variable map after all rotations have been applied.
//   - Log:  a slice of RotationEntry values recording what changed, including
//     the old value, the new value, and the timestamp of the rotation.
//
// Options allow callers to control edge-case behaviour:
//   - FailOnMissing: return an error when a rotation targets a key that does
//     not exist in the source map (default: false — the key is created).
//   - SkipUnchanged: omit log entries where the new value equals the existing
//     value (default: true).
//
// Typical usage:
//
//	result, err := envrotator.Rotate(vars, []envrotator.Rotation{
//	    {Key: "DB_PASSWORD", NewValue: newPassword},
//	}, envrotator.DefaultOptions())
package envrotator
