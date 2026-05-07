// Package envmap provides utility types and functions for working with
// environment variable collections as ordered Entry slices.
//
// It bridges the gap between the raw map[string]string representation used
// throughout envlayer and ordered, filterable, transformable sequences of
// key-value pairs.
//
// Core types:
//
//	- Entry: a single key/value pair
//
// Core functions:
//
//	- FromMap: convert a map to a sorted []Entry
//	- ToMap: convert []Entry back to a map
//	- Filter: select entries by key predicate
//	- MapValues: transform entry values with a function
package envmap
