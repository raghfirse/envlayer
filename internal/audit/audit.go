// Package audit provides change auditing for environment variable sets,
// tracking which keys were added, removed, or modified across env layers.
package audit

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// Entry records a single audited change to an environment variable.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Key       string    `json:"key"`
	Action    string    `json:"action"` // "added", "removed", "changed"
	OldValue  string    `json:"old_value,omitempty"`
	NewValue  string    `json:"new_value,omitempty"`
	Source    string    `json:"source,omitempty"`
}

// Log holds a sequence of audit entries.
type Log struct {
	Entries []Entry
}

// Record compares two env maps and appends audit entries for any differences.
// source identifies where the new values came from (e.g. a filename).
func (l *Log) Record(before, after map[string]string, source string) {
	now := time.Now().UTC()

	allKeys := unionKeys(before, after)
	for _, k := range allKeys {
		oldVal, inBefore := before[k]
		newVal, inAfter := after[k]

		switch {
		case inBefore && inAfter && oldVal != newVal:
			l.Entries = append(l.Entries, Entry{
				Timestamp: now,
				Key:       k,
				Action:    "changed",
				OldValue:  oldVal,
				NewValue:  newVal,
				Source:    source,
			})
		case !inBefore && inAfter:
			l.Entries = append(l.Entries, Entry{
				Timestamp: now,
				Key:       k,
				Action:    "added",
				NewValue:  newVal,
				Source:    source,
			})
		case inBefore && !inAfter:
			l.Entries = append(l.Entries, Entry{
				Timestamp: now,
				Key:       k,
				Action:    "removed",
				OldValue:  oldVal,
				Source:    source,
			})
		}
	}
}

// Print writes a human-readable summary of the audit log to w.
func (l *Log) Print(w io.Writer) {
	for _, e := range l.Entries {
		switch e.Action {
		case "added":
			fmt.Fprintf(w, "[%s] + %s = %q (source: %s)\n", e.Timestamp.Format(time.RFC3339), e.Key, e.NewValue, e.Source)
		case "removed":
			fmt.Fprintf(w, "[%s] - %s (was %q, source: %s)\n", e.Timestamp.Format(time.RFC3339), e.Key, e.OldValue, e.Source)
		case "changed":
			fmt.Fprintf(w, "[%s] ~ %s: %q -> %q (source: %s)\n", e.Timestamp.Format(time.RFC3339), e.Key, e.OldValue, e.NewValue, e.Source)
		}
	}
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
