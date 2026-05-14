// Package envpatcher applies partial updates (patches) to an existing
// environment variable map, supporting set, delete, and rename operations.
package envpatcher

import "fmt"

// OpKind identifies the type of patch operation.
type OpKind string

const (
	OpSet    OpKind = "set"
	OpDelete OpKind = "delete"
	OpRename OpKind = "rename"
)

// Op represents a single patch operation.
type Op struct {
	Kind    OpKind
	Key     string
	Value   string // used by OpSet
	NewKey  string // used by OpRename
}

// Result holds the patched map and a log of applied changes.
type Result struct {
	Vars    map[string]string
	Applied []string
	Skipped []string
}

// Apply executes the given ops against a copy of vars and returns a Result.
// The original map is never mutated.
func Apply(vars map[string]string, ops []Op) (*Result, error) {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		out[k] = v
	}

	result := &Result{Vars: out}

	for _, op := range ops {
		switch op.Kind {
		case OpSet:
			if op.Key == "" {
				return nil, fmt.Errorf("set op missing key")
			}
			out[op.Key] = op.Value
			result.Applied = append(result.Applied, fmt.Sprintf("set %s", op.Key))

		case OpDelete:
			if op.Key == "" {
				return nil, fmt.Errorf("delete op missing key")
			}
			if _, ok := out[op.Key]; !ok {
				result.Skipped = append(result.Skipped, fmt.Sprintf("delete %s (not found)", op.Key))
				continue
			}
			delete(out, op.Key)
			result.Applied = append(result.Applied, fmt.Sprintf("delete %s", op.Key))

		case OpRename:
			if op.Key == "" || op.NewKey == "" {
				return nil, fmt.Errorf("rename op requires both key and new_key")
			}
			val, ok := out[op.Key]
			if !ok {
				result.Skipped = append(result.Skipped, fmt.Sprintf("rename %s (not found)", op.Key))
				continue
			}
			out[op.NewKey] = val
			delete(out, op.Key)
			result.Applied = append(result.Applied, fmt.Sprintf("rename %s -> %s", op.Key, op.NewKey))

		default:
			return nil, fmt.Errorf("unknown op kind: %q", op.Kind)
		}
	}

	return result, nil
}
