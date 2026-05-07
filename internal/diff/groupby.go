package diff

import "sort"

// GroupByPrefix groups diff result keys by their prefix (e.g. "DB_", "AWS_").
// Keys without an underscore are placed under the empty string group "".
type GroupedResult struct {
	Prefix  string
	Keys    []string
}

// GroupResult partitions the keys from a Result into named prefix buckets.
// Only keys appearing in MissingInRight, MissingInLeft, or Mismatched are
// considered; keys that match on both sides are omitted.
func GroupResult(r Result) []GroupedResult {
	seen := make(map[string][]string)

	add := func(key string) {
		prefix := extractPrefix(key)
		seen[prefix] = append(seen[prefix], key)
	}

	for _, k := range r.MissingInRight {
		add(k)
	}
	for _, k := range r.MissingInLeft {
		add(k)
	}
	for _, m := range r.Mismatched {
		add(m.Key)
	}

	// Deduplicate within each prefix bucket.
	for prefix, keys := range seen {
		seen[prefix] = deduplicateStrings(keys)
	}

	// Build sorted slice of GroupedResult.
	prefixes := make([]string, 0, len(seen))
	for p := range seen {
		prefixes = append(prefixes, p)
	}
	sort.Strings(prefixes)

	out := make([]GroupedResult, 0, len(prefixes))
	for _, p := range prefixes {
		out = append(out, GroupedResult{Prefix: p, Keys: seen[p]})
	}
	return out
}

// extractPrefix returns the portion of key up to and including the first "_".
// If no underscore is present, the empty string is returned.
func extractPrefix(key string) string {
	for i, ch := range key {
		if ch == '_' {
			return key[:i+1]
		}
	}
	return ""
}

func deduplicateStrings(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	sort.Strings(out)
	return out
}
