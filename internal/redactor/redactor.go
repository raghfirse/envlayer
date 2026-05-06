// Package redactor provides functionality for redacting sensitive environment
// variable values from output strings such as log lines or command output.
package redactor

import (
	"strings"

	"github.com/yourusername/envlayer/internal/masker"
)

// Redactor holds a set of sensitive values and replaces them in arbitrary text.
type Redactor struct {
	sensitiveValues []string
	maskChar        string
}

// New creates a Redactor from a map of environment variables, extracting values
// for keys that are considered sensitive. maskChar is used as the replacement
// string; pass an empty string to use the default "****".
func New(vars map[string]string, maskChar string) *Redactor {
	if maskChar == "" {
		maskChar = "****"
	}

	var sensitive []string
	for k, v := range vars {
		if v != "" && masker.IsSensitive(k) {
			sensitive = append(sensitive, v)
		}
	}

	return &Redactor{
		sensitiveValues: sensitive,
		maskChar:        maskChar,
	}
}

// Redact replaces any sensitive values found in text with the mask character.
func (r *Redactor) Redact(text string) string {
	for _, val := range r.sensitiveValues {
		text = strings.ReplaceAll(text, val, r.maskChar)
	}
	return text
}

// RedactLines applies Redact to each line of a multi-line string and returns
// the result joined by newlines.
func (r *Redactor) RedactLines(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = r.Redact(line)
	}
	return strings.Join(lines, "\n")
}
