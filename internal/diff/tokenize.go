package diff

import (
	"sort"
	"strings"
)

// TokenizeOptions controls how values are tokenized for similarity analysis.
type TokenizeOptions struct {
	Delimiter string
	Lowercase bool
}

// DefaultTokenizeOptions returns sensible defaults.
func DefaultTokenizeOptions() TokenizeOptions {
	return TokenizeOptions{
		Delimiter: ",",
		Lowercase: true,
	}
}

// TokenizeResult holds the token sets extracted from two env maps.
type TokenizeResult struct {
	Key        string
	LeftTokens  []string
	RightTokens []string
	OnlyInLeft  []string
	OnlyInRight []string
	Shared      []string
}

// TokenizeValues splits values in both env maps by delimiter and computes
// per-key token-level diffs for keys that exist in both maps.
func TokenizeValues(left, right map[string]string, opts TokenizeOptions) []TokenizeResult {
	var results []TokenizeResult

	for key, lval := range left {
		rval, ok := right[key]
		if !ok {
			continue
		}
		if lval == rval {
			continue
		}

		ltoks := splitTokens(lval, opts)
		rtoks := splitTokens(rval, opts)

		lset := toTokenSet(ltoks)
		rset := toTokenSet(rtoks)

		var onlyLeft, onlyRight, shared []string
		for t := range lset {
			if rset[t] {
				shared = append(shared, t)
			} else {
				onlyLeft = append(onlyLeft, t)
			}
		}
		for t := range rset {
			if !lset[t] {
				onlyRight = append(onlyRight, t)
			}
		}

		sort.Strings(onlyLeft)
		sort.Strings(onlyRight)
		sort.Strings(shared)

		results = append(results, TokenizeResult{
			Key:         key,
			LeftTokens:  ltoks,
			RightTokens: rtoks,
			OnlyInLeft:  onlyLeft,
			OnlyInRight: onlyRight,
			Shared:      shared,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})
	return results
}

func splitTokens(val string, opts TokenizeOptions) []string {
	parts := strings.Split(val, opts.Delimiter)
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if opts.Lowercase {
			p = strings.ToLower(p)
		}
		out = append(out, p)
	}
	return out
}

func toTokenSet(tokens []string) map[string]bool {
	s := make(map[string]bool, len(tokens))
	for _, t := range tokens {
		s[t] = true
	}
	return s
}
