// Package merger provides utilities for combining multiple environment
// variable maps into a single resolved map.
//
// In envlayer, environment variables are loaded from multiple .env files
// (e.g., .env, .env.production, .env.local) and merged in a defined
// priority order. The merger package implements this layering logic.
//
// Typical usage:
//
//	base, _ := loader.LoadFile(".env")
//	prod, _ := loader.LoadFile(".env.production")
//	local, _ := loader.LoadFile(".env.local")
//
//	resolved := merger.Merge([]map[string]string{base, prod, local})
//
// In the example above, values in .env.local take highest precedence,
// followed by .env.production, and finally .env as the base defaults.
package merger
