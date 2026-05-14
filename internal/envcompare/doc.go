// Package envcompare provides structured comparison between two environment
// variable maps.
//
// It identifies keys that are added, removed, changed, or identical between
// a "left" and "right" snapshot, making it useful for auditing environment
// drift across deployments, profiles, or time-based snapshots.
//
// Example usage:
//
//	left  := map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
//	right := map[string]string{"DB_HOST": "prod.db",   "PORT": "5432", "LOG_LEVEL": "info"}
//
//	result := envcompare.Compare(left, right)
//	fmt.Println(envcompare.Summary(result))
//	// Output: 1 added, 1 changed, 1 identical
package envcompare
