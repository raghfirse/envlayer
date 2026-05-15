// Package envboundary enforces key-level access boundaries across environment
// variable maps.
//
// # Overview
//
// Each key in an env map can be assigned one of three access levels:
//
//   - Public     – safe to expose to any layer or consumer
//   - Internal   – restricted to within a defined service boundary
//   - Restricted – must never leak into downstream or exported contexts
//
// A Policy maps key prefixes or exact key names to their Level. When multiple
// prefix rules match, the longest prefix wins. Exact key matches always take
// precedence over prefix matches.
//
// # Usage
//
//	policy := envboundary.Policy{
//		Rules: map[string]envboundary.Level{
//			"SECRET_": envboundary.Restricted,
//			"APP_":    envboundary.Public,
//		},
//	}
//
//	// Check for violations when only Public keys are allowed:
//	violations := envboundary.Check(vars, policy, envboundary.Public)
//
//	// Produce a filtered map containing only Public keys:
//	safe := envboundary.Filter(vars, policy, envboundary.Public)
package envboundary
