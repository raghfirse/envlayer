package loader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap holds key-value pairs parsed from a .env file.
type EnvMap map[string]string

// LoadFile reads and parses a .env file at the given path.
// It supports KEY=VALUE and KEY="VALUE" formats, and ignores
// blank lines and lines starting with '#'.
func LoadFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("loader: open %q: %w", path, err)
	}
	defer f.Close()

	env := make(EnvMap)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("loader: %q line %d: %w", path, lineNum, err)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("loader: scanning %q: %w", path, err)
	}

	return env, nil
}

// parseLine splits a single KEY=VALUE line into its components.
func parseLine(line string) (string, string, error) {
	idx := strings.IndexByte(line, '=')
	if idx < 1 {
		return "", "", fmt.Errorf("invalid line %q: expected KEY=VALUE", line)
	}

	key := strings.TrimSpace(line[:idx])
	raw := strings.TrimSpace(line[idx+1:])

	value := stripQuotes(raw)
	return key, value, nil
}

// stripQuotes removes surrounding single or double quotes from a value.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
