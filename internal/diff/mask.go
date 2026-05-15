package diff

import (
	"regexp"
	"strings"
)

// MaskOptions controls how values are masked in output.
type MaskOptions struct {
	// MaskChar is the character used to replace masked characters.
	MaskChar string
	// RevealChars is the number of characters to reveal at the start of a value.
	RevealChars int
	// Patterns is a list of key substrings that trigger masking.
	Patterns []string
	// Regex is an optional compiled pattern; if set, keys matching it are masked.
	Regex *regexp.Regexp
}

// DefaultMaskOptions returns sensible defaults for masking sensitive values.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		MaskChar:    "*",
		RevealChars: 0,
		Patterns:    []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PRIVATE", "CREDENTIAL"},
	}
}

// shouldMask returns true when the key matches any configured pattern or regex.
func shouldMask(key string, opts MaskOptions) bool {
	upper := strings.ToUpper(key)
	for _, p := range opts.Patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	if opts.Regex != nil && opts.Regex.MatchString(key) {
		return true
	}
	return false
}

// maskValue replaces characters in value according to opts.
func maskValue(value string, opts MaskOptions) string {
	if len(value) == 0 {
		return value
	}
	reveal := opts.RevealChars
	if reveal < 0 {
		reveal = 0
	}
	if reveal >= len(value) {
		return value
	}
	visible := value[:reveal]
	masked := strings.Repeat(opts.MaskChar, len(value)-reveal)
	return visible + masked
}

// MaskEnv returns a copy of env with sensitive values masked.
func MaskEnv(env map[string]string, opts MaskOptions) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if shouldMask(k, opts) {
			out[k] = maskValue(v, opts)
		} else {
			out[k] = v
		}
	}
	return out
}

// MaskEnvWithReport masks env and returns the list of keys that were masked.
func MaskEnvWithReport(env map[string]string, opts MaskOptions) (map[string]string, []string) {
	out := make(map[string]string, len(env))
	var masked []string
	for k, v := range env {
		if shouldMask(k, opts) {
			out[k] = maskValue(v, opts)
			masked = append(masked, k)
		} else {
			out[k] = v
		}
	}
	sortStrings(masked)
	return out, masked
}
