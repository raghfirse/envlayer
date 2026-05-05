package loader

import "fmt"

// LoadFiles loads and returns the merged key-value pairs from multiple
// .env files in order. Each subsequent file's values are layered on top
// of the previous ones, with later files taking precedence.
//
// This is a convenience wrapper around LoadFile intended for use with
// the list of paths produced by the resolver package.
func LoadFiles(paths []string) (map[string]string, error) {
	result := make(map[string]string)

	for _, path := range paths {
		pairs, err := LoadFile(path)
		if err != nil {
			return nil, fmt.Errorf("loader: failed to load %q: %w", path, err)
		}
		for k, v := range pairs {
			result[k] = v
		}
	}

	return result, nil
}
