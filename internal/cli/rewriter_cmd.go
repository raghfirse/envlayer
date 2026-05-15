package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envlayer/internal/envrewriter"
	"github.com/user/envlayer/internal/loader"
)

// RunRewriter loads env files, applies rewrite rules supplied via flags, and
// prints the resulting variables to w.
//
// ruleStrs is a slice of "target:find:replace" strings, e.g.
//
//	"value:localhost:db.internal"
func RunRewriter(dir, env string, files []string, ruleStrs []string, format string, w io.Writer) error {
	if len(files) == 0 {
		return fmt.Errorf("envrewriter: at least one --file is required")
	}
	if len(ruleStrs) == 0 {
		return fmt.Errorf("envrewriter: at least one --rule is required")
	}

	// Resolve absolute paths relative to dir.
	paths := make([]string, len(files))
	for i, f := range files {
		if strings.HasPrefix(f, "/") {
			paths[i] = f
		} else {
			paths[i] = dir + "/" + f
		}
	}

	vars, err := loader.LoadFiles(paths)
	if err != nil {
		return fmt.Errorf("envrewriter: %w", err)
	}

	rules, err := parseRewriteRules(ruleStrs)
	if err != nil {
		return err
	}

	opts := envrewriter.DefaultOptions()
	opts.Rules = rules

	out, err := envrewriter.Rewrite(vars, opts)
	if err != nil {
		return fmt.Errorf("envrewriter: %w", err)
	}

	switch format {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	default:
		for _, k := range sortedRewriterKeys(out) {
			fmt.Fprintf(w, "%s=%s\n", k, out[k])
		}
	}
	return nil
}

func parseRewriteRules(strs []string) ([]envrewriter.Rule, error) {
	rules := make([]envrewriter.Rule, 0, len(strs))
	for _, s := range strs {
		parts := strings.SplitN(s, ":", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("envrewriter: invalid rule %q — expected target:find:replace", s)
		}
		target := parts[0]
		if target != "key" && target != "value" && target != "both" {
			return nil, fmt.Errorf("envrewriter: unknown target %q — must be key, value, or both", target)
		}
		rules = append(rules, envrewriter.Rule{Target: target, Find: parts[1], Replace: parts[2]})
	}
	return rules, nil
}

func sortedRewriterKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	os.Stderr.WriteString("") // satisfy import; actual sort below
	sortStringsLocal(keys)
	return keys
}

func sortStringsLocal(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
