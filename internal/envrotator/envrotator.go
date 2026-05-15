// Package envrotator provides utilities for rotating environment variable
// values — replacing old values with new ones across a variable map while
// maintaining a rotation log for auditing purposes.
package envrotator

import (
	"errors"
	"fmt"
	"time"
)

// RotationEntry records a single key rotation event.
type RotationEntry struct {
	Key       string    `json:"key"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
	RotatedAt time.Time `json:"rotated_at"`
}

// Result holds the rotated variable map and the log of changes made.
type Result struct {
	Vars    map[string]string
	Log     []RotationEntry
}

// Rotation describes a single key-to-new-value mapping.
type Rotation struct {
	Key      string
	NewValue string
}

// Options controls the behaviour of Rotate.
type Options struct {
	// FailOnMissing causes Rotate to return an error if a rotation targets a
	// key that does not exist in the source map.
	FailOnMissing bool

	// SkipUnchanged omits entries from the log where the new value equals
	// the existing value.
	SkipUnchanged bool
}

// DefaultOptions returns a sensible default Options value.
func DefaultOptions() Options {
	return Options{
		FailOnMissing: false,
		SkipUnchanged: true,
	}
}

// Rotate applies the given rotations to vars and returns a Result containing
// the updated map and a log of every change that was applied.
// The original vars map is never mutated.
func Rotate(vars map[string]string, rotations []Rotation, opts Options) (Result, error) {
	if len(rotations) == 0 {
		return Result{Vars: copyMap(vars)}, nil
	}

	out := copyMap(vars)
	var log []RotationEntry
	now := time.Now().UTC()

	for _, r := range rotations {
		if r.Key == "" {
			return Result{}, errors.New("envrotator: rotation key must not be empty")
		}

		old, exists := out[r.Key]
		if !exists {
			if opts.FailOnMissing {
				return Result{}, fmt.Errorf("envrotator: key %q not found in vars", r.Key)
			}
			// Key does not exist — create it.
			old = ""
		}

		if opts.SkipUnchanged && old == r.NewValue {
			continue
		}

		out[r.Key] = r.NewValue
		log = append(log, RotationEntry{
			Key:       r.Key,
			OldValue:  old,
			NewValue:  r.NewValue,
			RotatedAt: now,
		})
	}

	return Result{Vars: out, Log: log}, nil
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
