package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/your-org/envlayer/internal/envproxy"
	"github.com/your-org/envlayer/internal/loader"
	"github.com/your-org/envlayer/internal/resolver"
)

// RunProxy loads env files for the given environment and resolves one or more
// keys through a proxy that falls back to the real OS environment.
//
// Flags:
//
//	--dir      directory containing .env files (default ".")
//	--env      environment name (e.g. "production")
//	--keys     comma-separated list of keys to resolve
//	--os-fallback  if set, missing keys fall back to os.Environ
func RunProxy(args []string, out io.Writer) error {
	dir := "."
	env := ""
	keys := ""
	osFallback := false

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--dir":
			i++
			if i < len(args) {
				dir = args[i]
			}
		case "--env":
			i++
			if i < len(args) {
				env = args[i]
			}
		case "--keys":
			i++
			if i < len(args) {
				keys = args[i]
			}
		case "--os-fallback":
			osFallback = true
		}
	}

	if keys == "" {
		return fmt.Errorf("proxy: --keys is required")
	}

	files, err := resolver.Resolve(dir, env)
	if err != nil {
		return fmt.Errorf("proxy: %w", err)
	}

	vars, err := loader.LoadFiles(files)
	if err != nil {
		return fmt.Errorf("proxy: %w", err)
	}

	var p *envproxy.Proxy
	if osFallback {
		p = envproxy.WithOSFallback(vars)
	} else {
		p = envproxy.New(vars, nil)
	}

	for _, k := range strings.Split(keys, ",") {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		v, err := p.Get(k)
		if err != nil {
			_, _ = fmt.Fprintf(out, "%s=<not found>\n", k)
			continue
		}
		_, _ = fmt.Fprintf(out, "%s=%s\n", k, v)
	}

	return nil
}

// proxyOSFallback is a thin wrapper used in tests to inject a custom env.
func proxyOSFallback(environ []string) envproxy.FallbackFunc {
	m := make(map[string]string, len(environ))
	for _, e := range environ {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return func(key string) (string, bool) {
		v, ok := m[key]
		return v, ok
	}
}

// ensure proxyOSFallback is used (suppress unused warning in non-test builds).
var _ = proxyOSFallback(os.Environ())
