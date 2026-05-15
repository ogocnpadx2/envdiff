package diff

import (
	"fmt"
	"sort"
	"strings"
)

// ProfileEntry describes the inferred type/shape of a single key's value.
type ProfileEntry struct {
	Key      string
	Type     string // "bool", "int", "float", "url", "path", "empty", "string"
	Length   int
	NonEmpty bool
}

// ProfileResult holds the profiled entries for an env map.
type ProfileResult struct {
	Entries []ProfileEntry
}

// ProfileEnv analyses the types and shapes of values in an env map.
func ProfileEnv(env map[string]string) ProfileResult {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]ProfileEntry, 0, len(keys))
	for _, k := range keys {
		v := env[k]
		entries = append(entries, ProfileEntry{
			Key:      k,
			Type:     inferType(v),
			Length:   len(v),
			NonEmpty: v != "",
		})
	}
	return ProfileResult{Entries: entries}
}

func inferType(v string) string {
	if v == "" {
		return "empty"
	}
	lower := strings.ToLower(v)
	if lower == "true" || lower == "false" {
		return "bool"
	}
	if isInteger(v) {
		return "int"
	}
	if isFloat(v) {
		return "float"
	}
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
		return "url"
	}
	if strings.HasPrefix(v, "/") || strings.HasPrefix(v, "./") || strings.HasPrefix(v, "../") {
		return "path"
	}
	return "string"
}

func isInteger(s string) bool {
	if s == "" {
		return false
	}
	start := 0
	if s[0] == '-' || s[0] == '+' {
		start = 1
	}
	if start == len(s) {
		return false
	}
	for _, c := range s[start:] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func isFloat(s string) bool {
	if s == "" {
		return false
	}
	_, err := fmt.Sscanf(s, "%f", new(float64))
	return err == nil && strings.ContainsAny(s, ".eE")
}
