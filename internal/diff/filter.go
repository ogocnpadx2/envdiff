package diff

import "strings"

// FilterOptions controls which results are included in the output.
type FilterOptions struct {
	// OnlyMissing restricts output to keys missing in either file.
	OnlyMissing bool
	// OnlyMismatched restricts output to keys present in both files but with different values.
	OnlyMismatched bool
	// KeyPrefix filters results to keys that start with the given prefix (case-insensitive).
	KeyPrefix string
}

// FilterResult returns a new Result containing only the entries that match
// the provided FilterOptions. The original Result is not modified.
func FilterResult(r Result, opts FilterOptions) Result {
	out := Result{}

	if !opts.OnlyMismatched {
		for _, k := range r.MissingInRight {
			if matchesPrefix(k, opts.KeyPrefix) {
				out.MissingInRight = append(out.MissingInRight, k)
			}
		}
		for _, k := range r.MissingInLeft {
			if matchesPrefix(k, opts.KeyPrefix) {
				out.MissingInLeft = append(out.MissingInLeft, k)
			}
		}
	}

	if !opts.OnlyMissing {
		for _, m := range r.Mismatched {
			if matchesPrefix(m.Key, opts.KeyPrefix) {
				out.Mismatched = append(out.Mismatched, m)
			}
		}
	}

	return out
}

func matchesPrefix(key, prefix string) bool {
	if prefix == "" {
		return true
	}
	return strings.HasPrefix(strings.ToUpper(key), strings.ToUpper(prefix))
}
