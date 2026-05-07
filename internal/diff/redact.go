package diff

import (
	"regexp"
	"strings"
)

// RedactOptions controls which keys are redacted in output.
type RedactOptions struct {
	// Patterns is a list of substring or regex patterns to match key names.
	Patterns []string
	// Placeholder is substituted for sensitive values.
	Placeholder string
}

// DefaultRedactPatterns contains common sensitive key substrings.
var DefaultRedactPatterns = []string{
	"PASSWORD", "SECRET", "TOKEN", "API_KEY", "PRIVATE_KEY", "CREDENTIAL",
}

// DefaultPlaceholder is used when no custom placeholder is set.
const DefaultPlaceholder = "[REDACTED]"

// ParseRedactOptions builds a RedactOptions from CLI flag strings.
// patterns is a comma-separated list; placeholder defaults to [REDACTED].
func ParseRedactOptions(patterns, placeholder string) RedactOptions {
	var pats []string
	if patterns == "" {
		pats = append([]string{}, DefaultRedactPatterns...)
	} else {
		for _, p := range strings.Split(patterns, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				pats = append(pats, strings.ToUpper(p))
			}
		}
	}
	ph := placeholder
	if ph == "" {
		ph = DefaultPlaceholder
	}
	return RedactOptions{Patterns: pats, Placeholder: ph}
}

// shouldRedact returns true if the key matches any pattern.
func shouldRedact(key string, opts RedactOptions) bool {
	upper := strings.ToUpper(key)
	for _, p := range opts.Patterns {
		matched, err := regexp.MatchString(p, upper)
		if err == nil && matched {
			return true
		}
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// RedactResult returns a copy of Result with sensitive values replaced.
func RedactResult(r Result, opts RedactOptions) Result {
	out := Result{}

	for _, m := range r.Mismatched {
		entry := m
		if shouldRedact(m.Key, opts) {
			entry.LeftVal = opts.Placeholder
			entry.RightVal = opts.Placeholder
		}
		out.Mismatched = append(out.Mismatched, entry)
	}
	out.MissingInLeft = append([]string{}, r.MissingInLeft...)
	out.MissingInRight = append([]string{}, r.MissingInRight...)
	return out
}
