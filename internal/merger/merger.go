package merger

// Merger handles combining multiple environment variable maps
// with a defined priority order (later maps override earlier ones).

// Merge takes a slice of env maps and merges them in order.
// Keys in later maps override keys in earlier maps.
func Merge(layers []map[string]string) map[string]string {
	result := make(map[string]string)
	for _, layer := range layers {
		for k, v := range layer {
			result[k] = v
		}
	}
	return result
}

// MergeWithBase merges override on top of base, returning a new map.
// The base map is not modified.
func MergeWithBase(base, override map[string]string) map[string]string {
	result := make(map[string]string, len(base))
	for k, v := range base {
		result[k] = v
	}
	for k, v := range override {
		result[k] = v
	}
	return result
}

// Keys returns the sorted list of keys present in the merged map.
// Useful for deterministic output.
func Keys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sortStrings(keys)
	return keys
}

// sortStrings performs an in-place insertion sort on a string slice.
func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		key := s[i]
		j := i - 1
		for j >= 0 && s[j] > key {
			s[j+1] = s[j]
			j--
		}
		s[j+1] = key
	}
}
