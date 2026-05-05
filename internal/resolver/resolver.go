// Package resolver handles resolving the correct set of .env files
// for a given environment context (e.g., "development", "production").
package resolver

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the configuration for the resolver.
type Config struct {
	// BaseDir is the directory to search for .env files.
	BaseDir string
	// Environment is the target environment (e.g., "development", "staging").
	Environment string
}

// Resolve returns an ordered list of .env file paths that should be loaded
// and merged for the given environment. Files are returned in priority order:
// base .env first, then environment-specific overrides last.
func Resolve(cfg Config) ([]string, error) {
	if cfg.BaseDir == "" {
		cfg.BaseDir = "."
	}

	candidates := candidateFiles(cfg.BaseDir, cfg.Environment)

	var paths []string
	for _, c := range candidates {
		if fileExists(c) {
			paths = append(paths, c)
		}
	}

	if len(paths) == 0 {
		return nil, fmt.Errorf("resolver: no .env files found in %q for environment %q", cfg.BaseDir, cfg.Environment)
	}

	return paths, nil
}

// candidateFiles returns the ordered list of candidate file paths to check.
func candidateFiles(baseDir, env string) []string {
	candidates := []string{
		filepath.Join(baseDir, ".env"),
	}

	if env != "" {
		candidates = append(candidates,
			filepath.Join(baseDir, fmt.Sprintf(".env.%s", env)),
			filepath.Join(baseDir, fmt.Sprintf(".env.%s.local", env)),
		)
	}

	candidates = append(candidates, filepath.Join(baseDir, ".env.local"))

	return candidates
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
