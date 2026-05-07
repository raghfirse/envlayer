package pipeline

import (
	"strings"

	"github.com/user/envlayer/internal/interpolator"
	"github.com/user/envlayer/internal/masker"
	"github.com/user/envlayer/internal/transformer"
)

// StageInterpolate returns a Stage that expands variable references within
// values using the supplied base map as a fallback source.
func StageInterpolate(base map[string]string) Stage {
	return Stage{
		Name: "interpolate",
		Apply: func(vars map[string]string) (map[string]string, error) {
			return interpolator.Interpolate(vars, base), nil
		},
	}
}

// StageMask returns a Stage that replaces sensitive variable values with a
// masked placeholder using the default masker settings.
func StageMask(sensitiveKeys []string, maskChar string) Stage {
	opts := masker.Options{}
	if len(sensitiveKeys) > 0 {
		opts.ExtraKeys = sensitiveKeys
	}
	if maskChar != "" {
		opts.MaskChar = maskChar
	}
	return Stage{
		Name: "mask",
		Apply: func(vars map[string]string) (map[string]string, error) {
			return masker.Mask(vars, opts), nil
		},
	}
}

// StageAddPrefix returns a Stage that prepends prefix to every key.
func StageAddPrefix(prefix string) Stage {
	return Stage{
		Name: "add-prefix",
		Apply: func(vars map[string]string) (map[string]string, error) {
			return transformer.Transform(vars, transformer.Options{AddPrefix: prefix}), nil
		},
	}
}

// StageUppercaseKeys returns a Stage that uppercases all keys.
func StageUppercaseKeys() Stage {
	return Stage{
		Name: "uppercase-keys",
		Apply: func(vars map[string]string) (map[string]string, error) {
			return transformer.Transform(vars, transformer.Options{UppercaseKeys: true}), nil
		},
	}
}

// StageFilterPrefix returns a Stage that keeps only keys that start with the
// given prefix, optionally stripping the prefix from the resulting keys.
func StageFilterPrefix(prefix string, strip bool) Stage {
	return Stage{
		Name: "filter-prefix",
		Apply: func(vars map[string]string) (map[string]string, error) {
			out := make(map[string]string)
			for k, v := range vars {
				if strings.HasPrefix(k, prefix) {
					key := k
					if strip {
						key = strings.TrimPrefix(k, prefix)
					}
					out[key] = v
				}
			}
			return out, nil
		},
	}
}
