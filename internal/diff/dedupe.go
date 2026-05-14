package diff

import "sort"

// DedupeOptions controls how duplicate key detection works.
type DedupeOptions struct {
	// CaseSensitive determines whether key comparison is case-sensitive.
	CaseSensitive bool
}

// DuplicateEntry represents a key that appears more than once in an env map.
type DuplicateEntry struct {
	Key    string
	Values []string
}

// DedupeResult holds the outcome of a deduplication pass.
type DedupeResult struct {
	Duplicates []DuplicateEntry
	Clean      map[string]string
}

// DefaultDedupeOptions returns sensible defaults.
func DefaultDedupeOptions() DedupeOptions {
	return DedupeOptions{CaseSensitive: true}
}

// FindDuplicates scans a slice of raw key=value lines and identifies keys
// declared more than once, returning both the duplicates and a clean map
// that retains the last-seen value for each key (matching shell semantics).
func FindDuplicates(lines []string, opts DedupeOptions) DedupeResult {
	seen := map[string][]string{} // normalised key -> all raw values
	normKey := func(k string) string {
		if opts.CaseSensitive {
			return k
		}
		return toLower(k)
	}

	for _, line := range lines {
		key, val, ok := splitLine(line)
		if !ok {
			continue
		}
		nk := normKey(key)
		seen[nk] = append(seen[nk], val)
	}

	var dups []DuplicateEntry
	clean := map[string]string{}

	for nk, vals := range seen {
		clean[nk] = vals[len(vals)-1]
		if len(vals) > 1 {
			dups = append(dups, DuplicateEntry{Key: nk, Values: vals})
		}
	}

	sort.Slice(dups, func(i, j int) bool {
		return dups[i].Key < dups[j].Key
	})

	return DedupeResult{Duplicates: dups, Clean: clean}
}

// splitLine parses "KEY=VALUE" and returns (key, value, ok).
func splitLine(line string) (string, string, bool) {
	for i := 0; i < len(line); i++ {
		if line[i] == '=' {
			return line[:i], line[i+1:], true
		}
	}
	return "", "", false
}

// toLower is a simple ASCII lowercase helper.
func toLower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}
