// Package envexpander provides utilities for expanding abbreviated or
// shorthand environment variable keys into their full canonical names.
//
// An expansion map defines the mapping from short key to canonical key:
//
//	expansions := map[string]string{
//		"DB_HOST": "DATABASE_HOST",
//		"DB_PORT": "DATABASE_PORT",
//	}
//
// Use Expand to apply the map to a set of environment variables:
//
//	out, err := envexpander.Expand(vars, expansions, envexpander.DefaultOptions())
//
// Use Invert to reverse an expansion map, mapping canonical names back to
// their shorthand equivalents.
package envexpander
