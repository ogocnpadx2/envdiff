package diff

import (
	"bufio"
	"os"
	"strings"
)

// IgnoreOptions controls which keys are excluded from diff results.
type IgnoreOptions struct {
	// Keys is an explicit list of key names to ignore.
	Keys []string
	// Prefixes is a list of key prefixes; any key matching one is ignored.
	Prefixes []string
}

// ParseIgnoreFile reads a file where each non-blank, non-comment line is a key
// or prefix pattern (prefix patterns end with "*") to ignore.
func ParseIgnoreFile(path string) (IgnoreOptions, error) {
	f, err := os.Open(path)
	if err != nil {
		return IgnoreOptions{}, err
	}
	defer f.Close()

	var opts IgnoreOptions
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasSuffix(line, "*") {
			opts.Prefixes = append(opts.Prefixes, strings.TrimSuffix(line, "*"))
		} else {
			opts.Keys = append(opts.Keys, line)
		}
	}
	return opts, scanner.Err()
}

// shouldIgnore returns true if key matches any ignore rule.
func shouldIgnore(key string, opts IgnoreOptions) bool {
	for _, k := range opts.Keys {
		if k == key {
			return true
		}
	}
	for _, p := range opts.Prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}

// ApplyIgnore removes ignored keys from all slices in a Result.
func ApplyIgnore(r Result, opts IgnoreOptions) Result {
	if len(opts.Keys) == 0 && len(opts.Prefixes) == 0 {
		return r
	}

	filtered := Result{}

	for _, k := range r.MissingInRight {
		if !shouldIgnore(k, opts) {
			filtered.MissingInRight = append(filtered.MissingInRight, k)
		}
	}
	for _, k := range r.MissingInLeft {
		if !shouldIgnore(k, opts) {
			filtered.MissingInLeft = append(filtered.MissingInLeft, k)
		}
	}
	for _, m := range r.Mismatched {
		if !shouldIgnore(m.Key, opts) {
			filtered.Mismatched = append(filtered.Mismatched, m)
		}
	}
	return filtered
}
