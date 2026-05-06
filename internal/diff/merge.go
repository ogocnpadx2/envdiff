package diff

import (
	"fmt"
	"sort"
)

// MergeStrategy defines how to resolve conflicts when merging env files.
type MergeStrategy int

const (
	// MergePreferLeft keeps the left-side value on conflict.
	MergePreferLeft MergeStrategy = iota
	// MergePreferRight keeps the right-side value on conflict.
	MergePreferRight
	// MergeUnionAll includes all keys from both sides.
	MergeUnionAll
)

// ParseMergeStrategy parses a string into a MergeStrategy.
func ParseMergeStrategy(s string) (MergeStrategy, error) {
	switch s {
	case "left":
		return MergePreferLeft, nil
	case "right":
		return MergePreferRight, nil
	case "union":
		return MergeUnionAll, nil
	default:
		return 0, fmt.Errorf("unknown merge strategy %q: must be left, right, or union", s)
	}
}

// Merge combines two env maps into a single map using the given strategy.
// The returned map is a new map and does not modify the inputs.
func Merge(left, right map[string]string, strategy MergeStrategy) map[string]string {
	result := make(map[string]string)

	for k, v := range left {
		result[k] = v
	}

	for k, v := range right {
		if _, exists := result[k]; !exists {
			result[k] = v
			continue
		}
		// Key exists in both — apply strategy
		switch strategy {
		case MergePreferRight:
			result[k] = v
		case MergePreferLeft:
			// already set from left, keep it
		case MergeUnionAll:
			// keep left value but ensure key is present (already set)
		}
	}

	return result
}

// MergedKeys returns a sorted slice of all keys in the merged map.
func MergedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
