// Package interpolator provides variable interpolation for environment maps.
// It resolves references like ${VAR} or $VAR within values using other
// keys in the same map or a provided base environment.
package interpolator

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var varPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Options controls interpolation behaviour.
type Options struct {
	// FallbackToOS allows unresolved references to be looked up in os.Environ.
	FallbackToOS bool
	// ErrorOnMissing returns an error when a referenced variable cannot be resolved.
	ErrorOnMissing bool
}

// Interpolate resolves variable references in all values of env.
// Resolution order: env itself, then base (if provided), then OS (if opted in).
func Interpolate(env map[string]string, base map[string]string, opts Options) (map[string]string, error) {
	result := make(map[string]string, len(env))
	for k, v := range env {
		resolved, err := resolve(v, env, base, opts)
		if err != nil {
			return nil, fmt.Errorf("interpolator: key %q: %w", k, err)
		}
		result[k] = resolved
	}
	return result, nil
}

func resolve(value string, env, base map[string]string, opts Options) (string, error) {
	var resolveErr error
	result := varPattern.ReplaceAllStringFunc(value, func(match string) string {
		if resolveErr != nil {
			return match
		}
		name := extractName(match)
		if v, ok := env[name]; ok {
			return v
		}
		if base != nil {
			if v, ok := base[name]; ok {
				return v
			}
		}
		if opts.FallbackToOS {
			if v, ok := os.LookupEnv(name); ok {
				return v
			}
		}
		if opts.ErrorOnMissing {
			resolveErr = fmt.Errorf("undefined variable %q", name)
			return match
		}
		return ""
	})
	if resolveErr != nil {
		return "", resolveErr
	}
	return result, nil
}

func extractName(match string) string {
	match = strings.TrimPrefix(match, "${") 
	match = strings.TrimSuffix(match, "}")
	match = strings.TrimPrefix(match, "$")
	return match
}
