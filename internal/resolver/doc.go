// Package resolver provides environment-aware .env file path resolution
// for the envlayer tool.
//
// Given a base directory and an environment name (e.g., "production"),
// the resolver determines which .env files exist and returns them in
// the correct merge-priority order:
//
//  1. .env          — shared base variables
//  2. .env.<env>    — environment-specific overrides
//  3. .env.<env>.local — local machine overrides for the environment
//  4. .env.local    — local machine overrides (all environments)
//
// Files that do not exist on disk are silently skipped. If no files
// are found at all, Resolve returns an error.
package resolver
