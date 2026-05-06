// Package snapshot provides utilities for capturing and persisting
// the merged state of environment variables at a point in time.
//
// A snapshot records the full set of resolved environment variables
// along with metadata such as the target environment name and
// creation timestamp. Snapshots can be saved to disk as JSON files
// and reloaded later for auditing, diffing, or reproducibility.
//
// Typical usage:
//
//	snap := snapshot.Take(mergedEnv, "production")
//	if err := snapshot.Save(snap, ".env.snapshot.json"); err != nil {
//		log.Fatal(err)
//	}
//
// To reload and inspect a previous snapshot:
//
//	snap, err := snapshot.Load(".env.snapshot.json")
package snapshot
