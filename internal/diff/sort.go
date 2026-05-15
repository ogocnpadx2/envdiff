package diff

import (
	"fmt"
	"sort"
	"strings"
)

// SortField represents a field by which env entries can be sorted.
type SortField string

const (
	SortByKey    SortField = "key"
	SortByValue  SortField = "value"
	SortByLength SortField = "length"
)

// SortOrder represents ascending or descending order.
type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

// SortOptions configures how SortEnv behaves.
type SortOptions struct {
	Field     SortField
	Order     SortOrder
	IgnoreCase bool
}

// DefaultSortOptions returns sensible defaults.
func DefaultSortOptions() SortOptions {
	return SortOptions{
		Field:      SortByKey,
		Order:      SortAsc,
		IgnoreCase: false,
	}
}

// ParseSortField parses a field string into a SortField.
func ParseSortField(s string) (SortField, error) {
	switch strings.ToLower(s) {
	case "key":
		return SortByKey, nil
	case "value":
		return SortByValue, nil
	case "length":
		return SortByLength, nil
	default:
		return "", fmt.Errorf("unknown sort field %q: must be key, value, or length", s)
	}
}

// ParseSortOrder parses an order string into a SortOrder.
func ParseSortOrder(s string) (SortOrder, error) {
	switch strings.ToLower(s) {
	case "asc", "":
		return SortAsc, nil
	case "desc":
		return SortDesc, nil
	default:
		return "", fmt.Errorf("unknown sort order %q: must be asc or desc", s)
	}
}

// SortEnvEntry holds a key-value pair for sorted output.
type SortEnvEntry struct {
	Key   string
	Value string
}

// SortEnv returns env map entries sorted according to opts.
func SortEnv(env map[string]string, opts SortOptions) []SortEnvEntry {
	entries := make([]SortEnvEntry, 0, len(env))
	for k, v := range env {
		entries = append(entries, SortEnvEntry{Key: k, Value: v})
	}

	sort.SliceStable(entries, func(i, j int) bool {
		var less bool
		switch opts.Field {
		case SortByValue:
			li, lj := entries[i].Value, entries[j].Value
			if opts.IgnoreCase {
				li, lj = strings.ToLower(li), strings.ToLower(lj)
			}
			less = li < lj
		case SortByLength:
			less = len(entries[i].Value) < len(entries[j].Value)
		default: // SortByKey
			li, lj := entries[i].Key, entries[j].Key
			if opts.IgnoreCase {
				li, lj = strings.ToLower(li), strings.ToLower(lj)
			}
			less = li < lj
		}
		if opts.Order == SortDesc {
			return !less
		}
		return less
	})

	return entries
}
