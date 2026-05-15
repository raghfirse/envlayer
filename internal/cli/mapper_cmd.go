package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/your-org/envlayer/internal/envmapper"
	"github.com/your-org/envlayer/internal/loader"
)

// RunMapper loads env files, applies key-remapping rules, and prints the
// result in dotenv format.
//
// ruleArgs is a slice of "FROM:TO" or "FROM:TO:keep" strings.
func RunMapper(dir, env string, files []string, ruleArgs []string, dropUnmapped, failOnMissing bool, w io.Writer) error {
	if len(files) == 0 {
		return fmt.Errorf("mapper: at least one env file is required")
	}
	if len(ruleArgs) == 0 {
		return fmt.Errorf("mapper: at least one mapping rule is required")
	}

	rules, err := parseMapperRules(ruleArgs)
	if err != nil {
		return err
	}

	var paths []string
	for _, f := range files {
		paths = append(paths, dir+"/"+f)
	}

	vars, err := loader.LoadFiles(paths)
	if err != nil {
		return fmt.Errorf("mapper: %w", err)
	}

	opts := envmapper.Options{
		Rules:         rules,
		DropUnmapped:  dropUnmapped,
		FailOnMissing: failOnMissing,
	}

	out, err := envmapper.Apply(vars, opts)
	if err != nil {
		return fmt.Errorf("mapper: %w", err)
	}

	keys := envmapper.Keys(opts)
	if !dropUnmapped {
		// also print non-rule keys
		seen := make(map[string]bool)
		for _, k := range keys {
			seen[k] = true
		}
		for k := range out {
			if !seen[k] {
				keys = append(keys, k)
			}
		}
	}

	for _, k := range keys {
		v, ok := out[k]
		if !ok {
			continue
		}
		fmt.Fprintf(w, "%s=%s\n", k, v)
	}
	return nil
}

func parseMapperRules(args []string) ([]envmapper.Rule, error) {
	rules := make([]envmapper.Rule, 0, len(args))
	for _, a := range args {
		parts := strings.SplitN(a, ":", 3)
		if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("mapper: invalid rule %q — expected FROM:TO[:keep]", a)
		}
		rule := envmapper.Rule{From: parts[0], To: parts[1]}
		if len(parts) == 3 && strings.EqualFold(parts[2], "keep") {
			rule.Keep = true
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

// Ensure os is used (for potential future exit-code helpers).
var _ = os.Stderr
