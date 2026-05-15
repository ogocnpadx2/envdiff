package diff

import "strings"

// TruncateOptions controls how values are truncated in output.
type TruncateOptions struct {
	// MaxLength is the maximum number of characters to show per value.
	MaxLength int
	// Suffix is appended when a value is truncated (default: "...").
	Suffix string
	// KeyFilter limits truncation to keys containing this substring (empty = all keys).
	KeyFilter string
}

// DefaultTruncateOptions returns sensible defaults.
func DefaultTruncateOptions() TruncateOptions {
	return TruncateOptions{
		MaxLength: 40,
		Suffix:    "...",
		KeyFilter: "",
	}
}

// TruncateEnv returns a copy of env with values truncated according to opts.
func TruncateEnv(env map[string]string, opts TruncateOptions) map[string]string {
	if opts.MaxLength <= 0 {
		opts.MaxLength = DefaultTruncateOptions().MaxLength
	}
	if opts.Suffix == "" && opts.MaxLength > 0 {
		opts.Suffix = "..."
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		if opts.KeyFilter != "" && !strings.Contains(k, opts.KeyFilter) {
			out[k] = v
			continue
		}
		if len(v) > opts.MaxLength {
			out[k] = v[:opts.MaxLength] + opts.Suffix
		} else {
			out[k] = v
		}
	}
	return out
}

// TruncateEnvWithReport returns the truncated env and a list of keys that were truncated.
func TruncateEnvWithReport(env map[string]string, opts TruncateOptions) (map[string]string, []string) {
	if opts.MaxLength <= 0 {
		opts.MaxLength = DefaultTruncateOptions().MaxLength
	}
	if opts.Suffix == "" {
		opts.Suffix = "..."
	}

	out := make(map[string]string, len(env))
	var truncated []string

	for k, v := range env {
		if opts.KeyFilter != "" && !strings.Contains(k, opts.KeyFilter) {
			out[k] = v
			continue
		}
		if len(v) > opts.MaxLength {
			out[k] = v[:opts.MaxLength] + opts.Suffix
			truncated = append(truncated, k)
		} else {
			out[k] = v
		}
	}
	sortStrings(truncated)
	return out, truncated
}
