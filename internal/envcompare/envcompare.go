// Package envcompare provides utilities for comparing two env maps
// across multiple dimensions: keys, values, types, and structure.
package envcompare

import "sort"

// Result holds the full comparison between two env maps.
type Result struct {
	OnlyInLeft  []string          // keys present only in left
	OnlyInRight []string          // keys present only in right
	InBoth      []string          // keys present in both
	Changed     map[string][2]string // key -> [leftVal, rightVal] for changed keys
	Identical   []string          // keys with identical values
}

// Compare performs a full structural comparison of two env maps.
func Compare(left, right map[string]string) Result {
	r := Result{
		Changed:   make(map[string][2]string),
	}

	leftSet := make(map[string]bool, len(left))
	for k := range left {
		leftSet[k] = true
	}

	rightSet := make(map[string]bool, len(right))
	for k := range right {
		rightSet[k] = true
	}

	for k := range leftSet {
		if !rightSet[k] {
			r.OnlyInLeft = append(r.OnlyInLeft, k)
		} else {
			r.InBoth = append(r.InBoth, k)
			if left[k] != right[k] {
				r.Changed[k] = [2]string{left[k], right[k]}
			} else {
				r.Identical = append(r.Identical, k)
			}
		}
	}

	for k := range rightSet {
		if !leftSet[k] {
			r.OnlyInRight = append(r.OnlyInRight, k)
		}
	}

	sort.Strings(r.OnlyInLeft)
	sort.Strings(r.OnlyInRight)
	sort.Strings(r.InBoth)
	sort.Strings(r.Identical)
	return r
}

// Equal returns true if both maps contain exactly the same keys and values.
func Equal(left, right map[string]string) bool {
	if len(left) != len(right) {
		return false
	}
	for k, v := range left {
		if right[k] != v {
			return false
		}
	}
	return true
}

// Summary returns a human-readable summary string of the comparison result.
func Summary(r Result) string {
	added := len(r.OnlyInRight)
	removed := len(r.OnlyInLeft)
	changed := len(r.Changed)
	identical := len(r.Identical)

	if added == 0 && removed == 0 && changed == 0 {
		return "no differences found"
	}

	s := ""
	if added > 0 {
		s += sprint(added, "added")
	}
	if removed > 0 {
		if s != "" { s += ", " }
		s += sprint(removed, "removed")
	}
	if changed > 0 {
		if s != "" { s += ", " }
		s += sprint(changed, "changed")
	}
	if identical > 0 {
		if s != "" { s += ", " }
		s += sprint(identical, "identical")
	}
	return s
}

func sprint(n int, label string) string {
	if n == 1 {
		return "1 " + label
	}
	return itoa(n) + " " + label
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
