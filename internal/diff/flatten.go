package diff

import (
	"fmt"
	"sort"
	"strings"
)

// FlattenOptions controls how nested key structures are flattened.
type FlattenOptions struct {
	// Separator is the delimiter used to join nested key segments (default: "__").
	Separator string
	// MaxDepth limits how many segments are split (0 = unlimited).
	MaxDepth int
	// Uppercase converts all resulting keys to uppercase.
	Uppercase bool
}

// DefaultFlattenOptions returns sensible defaults.
func DefaultFlattenOptions() FlattenOptions {
	return FlattenOptions{
		Separator: "__",
		MaxDepth:  0,
		Uppercase: false,
	}
}

// FlattenResult holds the output of a flatten operation.
type FlattenResult struct {
	// Flattened is the resulting key→value map after flattening.
	Flattened map[string]string
	// Renamed tracks original→new key mappings where the key changed.
	Renamed map[string]string
	// Unchanged lists keys that were not affected.
	Unchanged []string
}

// FlattenEnv takes an env map and rewrites keys by collapsing repeated
// separator segments into a single canonical key. For example, with
// separator "__", "DB__HOST__PORT" stays as-is but "DB____HOST" becomes
// "DB__HOST" (consecutive separators are collapsed).
func FlattenEnv(env map[string]string, opts FlattenOptions) FlattenResult {
	if opts.Separator == "" {
		opts.Separator = "__"
	}

	result := FlattenResult{
		Flattened: make(map[string]string, len(env)),
		Renamed:   make(map[string]string),
	}

	for origKey, val := range env {
		newKey := normalizeKey(origKey, opts)
		result.Flattened[newKey] = val
		if newKey != origKey {
			result.Renamed[origKey] = newKey
		} else {
			result.Unchanged = append(result.Unchanged, origKey)
		}
	}

	sort.Strings(result.Unchanged)
	return result
}

// normalizeKey collapses consecutive separators and optionally uppercases.
func normalizeKey(key string, opts FlattenOptions) string {
	sep := opts.Separator
	segments := strings.Split(key, sep)

	// Remove empty segments produced by consecutive separators.
	filtered := segments[:0]
	for _, s := range segments {
		if s != "" {
			filtered = append(filtered, s)
		}
	}

	if opts.MaxDepth > 0 && len(filtered) > opts.MaxDepth {
		// Join the tail back into the last segment.
		tail := strings.Join(filtered[opts.MaxDepth-1:], sep)
		filtered = append(filtered[:opts.MaxDepth-1], tail)
	}

	out := strings.Join(filtered, sep)
	if opts.Uppercase {
		out = strings.ToUpper(out)
	}
	return out
}

// FormatFlattenText returns a human-readable summary of a FlattenResult.
func FormatFlattenText(r FlattenResult) string {
	var sb strings.Builder
	if len(r.Renamed) == 0 {
		sb.WriteString("No keys required flattening.\n")
		return sb.String()
	}

	// Stable output: sort by original key.
	origKeys := make([]string, 0, len(r.Renamed))
	for k := range r.Renamed {
		origKeys = append(origKeys, k)
	}
	sort.Strings(origKeys)

	fmt.Fprintf(&sb, "Flattened %d key(s):\n", len(r.Renamed))
	for _, orig := range origKeys {
		fmt.Fprintf(&sb, "  %s  →  %s\n", orig, r.Renamed[orig])
	}
	return sb.String()
}
