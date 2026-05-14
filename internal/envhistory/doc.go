// Package envhistory provides persistent history tracking for environment
// variable snapshots.
//
// Each recorded entry is stored as a JSON file in a configurable directory,
// capturing the full set of variables alongside a human-readable label and
// creation timestamp.
//
// Usage:
//
//	// Record the current state
//	e, err := envhistory.Record("/var/envlayer/history", "deploy-v1.2", vars)
//
//	// List all recorded entries (oldest first)
//	entries, err := envhistory.List("/var/envlayer/history")
//
//	// Retrieve a specific entry by ID
//	e, err := envhistory.Get("/var/envlayer/history", id)
//
//	// Remove an entry
//	err = envhistory.Delete("/var/envlayer/history", id)
package envhistory
