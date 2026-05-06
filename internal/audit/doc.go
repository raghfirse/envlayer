// Package audit records and reports changes between environment variable maps
// as they are merged across layers in envlayer.
//
// Usage:
//
//	l := &audit.Log{}
//	l.Record(baseEnv, mergedEnv, ".env.production")
//	l.Print(os.Stdout)
//
// Each call to Record compares the before and after maps and appends
// structured Entry values describing keys that were added, removed,
// or changed. Entries are ordered by key name within each call.
//
// The Print method writes a human-readable summary to any io.Writer,
// suitable for CLI output or log files.
package audit
