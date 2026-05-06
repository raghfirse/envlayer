// Package masker provides utilities for masking sensitive environment variable
// values before they are printed, exported, or logged.
//
// Sensitive keys are identified by substring matching against a configurable
// list of patterns (e.g. "PASSWORD", "SECRET", "TOKEN"). Matching is
// case-insensitive so both DB_PASSWORD and db_password are caught.
//
// Example usage:
//
//	safe := masker.Mask(vars, masker.Options{
//		MaskChar: "[hidden]",
//	})
//
// The original map is never mutated; Mask always returns a new map.
package masker
