// Package envrewriter provides find-and-replace rewriting for environment
// variable maps.
//
// Rules can target keys, values, or both, and support both case-sensitive
// and case-insensitive matching. Multiple rules are applied in order, so
// later rules operate on the output of earlier ones.
//
// Example:
//
//	opts := envrewriter.DefaultOptions()
//	opts.Rules = []envrewriter.Rule{
//		{Target: "key",   Find: "DEV_",  Replace: "PROD_"},
//		{Target: "value", Find: "localhost", Replace: "prod.db"},
//	}
//	out, err := envrewriter.Rewrite(vars, opts)
package envrewriter
