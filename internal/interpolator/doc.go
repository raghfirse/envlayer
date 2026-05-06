// Package interpolator resolves variable references embedded in environment
// variable values.
//
// It supports two reference syntaxes:
//
//	${VAR_NAME}   — braced form, recommended for clarity
//	$VAR_NAME     — unbraced form, resolved greedily
//
// Resolution order:
//  1. The current env map itself (peer variables).
//  2. An optional base map (e.g. a .env.base layer).
//  3. The host OS environment, when Options.FallbackToOS is true.
//
// Unresolved references are silently replaced with an empty string unless
// Options.ErrorOnMissing is set, in which case an error is returned.
package interpolator
