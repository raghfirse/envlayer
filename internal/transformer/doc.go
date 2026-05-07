// Package transformer applies structural transformations to environment
// variable maps.
//
// Supported transformations include:
//
//   - AddPrefix    – prepend a string to every key
//   - StripPrefix  – remove a leading string from every key
//   - UppercaseKeys – normalise keys to UPPER_CASE
//   - LowercaseKeys – normalise keys to lower_case
//   - TrimValues   – strip surrounding whitespace from values
//
// All transformations are non-destructive: the original map is never
// modified. A new map is always returned.
package transformer
