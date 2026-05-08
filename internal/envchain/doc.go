// Package envchain implements a priority-ordered chain of environment variable
// layers.
//
// A Chain holds multiple maps where index 0 is the highest-priority layer.
// Key lookups walk the chain from highest to lowest priority, returning the
// first match found — similar to prototype-chain resolution in JavaScript.
//
// Example usage:
//
//	base := map[string]string{"PORT": "8080", "LOG_LEVEL": "info"}
//	prod := map[string]string{"LOG_LEVEL": "warn"}
//
//	chain := envchain.New(prod, base)
//	chain.Get("LOG_LEVEL")   // → "warn"
//	chain.Get("PORT")        // → "8080"
//	chain.Resolve()          // → merged map, prod wins conflicts
package envchain
