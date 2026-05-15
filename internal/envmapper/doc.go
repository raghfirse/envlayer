// Package envmapper remaps environment variable keys according to a
// declarative set of rules.
//
// Each Rule specifies a source key (From) and a destination key (To).
// By default the source key is removed from the output map once remapped;
// setting Rule.Keep retains both the original and the new key.
//
// Options.DropUnmapped causes any key not referenced by a rule to be
// excluded from the output, producing a tightly scoped result map.
//
// Options.FailOnMissing makes Apply return an error when a rule's source
// key does not exist in the input, useful for strict pipeline validation.
//
// Example:
//
//	opts := envmapper.Options{
//		Rules: []envmapper.Rule{
//			{From: "DB_HOST", To: "DATABASE_HOST"},
//			{From: "DB_PASS", To: "DATABASE_PASSWORD", Keep: true},
//		},
//	}
//	out, err := envmapper.Apply(vars, opts)
package envmapper
