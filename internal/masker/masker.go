package masker

import "strings"

// DefaultSensitiveKeys contains common patterns for sensitive environment variables.
var DefaultSensitiveKeys = []string{
	"PASSWORD", "SECRET", "TOKEN", "KEY", "PRIVATE", "CREDENTIAL", "AUTH", "API_KEY",
}

// Options configures masking behaviour.
type Options struct {
	// SensitiveKeys is a list of substrings; any env var whose uppercased name
	// contains one of these substrings will have its value masked.
	SensitiveKeys []string
	// MaskChar is the string used to replace sensitive values. Defaults to "****".
	MaskChar string
}

// Mask returns a copy of vars where sensitive values are replaced with the
// mask character. Keys are matched case-insensitively against SensitiveKeys.
func Mask(vars map[string]string, opts Options) map[string]string {
	if opts.MaskChar == "" {
		opts.MaskChar = "****"
	}
	if len(opts.SensitiveKeys) == 0 {
		opts.SensitiveKeys = DefaultSensitiveKeys
	}

	result := make(map[string]string, len(vars))
	for k, v := range vars {
		if isSensitive(k, opts.SensitiveKeys) {
			result[k] = opts.MaskChar
		} else {
			result[k] = v
		}
	}
	return result
}

// IsSensitive reports whether the given key name is considered sensitive.
func IsSensitive(key string, sensitiveKeys []string) bool {
	return isSensitive(key, sensitiveKeys)
}

func isSensitive(key string, sensitiveKeys []string) bool {
	upper := strings.ToUpper(key)
	for _, s := range sensitiveKeys {
		if strings.Contains(upper, strings.ToUpper(s)) {
			return true
		}
	}
	return false
}
