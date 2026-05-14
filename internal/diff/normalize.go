package diff

import (
	"sort"
	"strings"
)

// NormalizeOptions controls how env map values are normalized before comparison.
type NormalizeOptions struct {
	TrimSpace    bool
	LowercaseVal bool
	LowercaseKey bool
	CollapseEmpty bool
}

// DefaultNormalizeOptions returns sensible defaults: trim whitespace only.
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		TrimSpace:    true,
		LowercaseVal: false,
		LowercaseKey: false,
		CollapseEmpty: false,
	}
}

// NormalizeEnv applies normalization rules to a parsed env map and returns
// a new map with the transformed keys/values.
func NormalizeEnv(env map[string]string, opts NormalizeOptions) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if opts.TrimSpace {
			v = strings.TrimSpace(v)
			k = strings.TrimSpace(k)
		}
		if opts.LowercaseKey {
			k = strings.ToLower(k)
		}
		if opts.LowercaseVal {
			v = strings.ToLower(v)
		}
		if opts.CollapseEmpty && v == "" {
			continue
		}
		out[k] = v
	}
	return out
}

// NormalizeReport summarises what changed during normalization.
type NormalizeReport struct {
	Trimmed   []string // keys whose values were trimmed
	Renamed   []string // keys that were lowercased
	Dropped   []string // keys dropped due to CollapseEmpty
}

// NormalizeEnvWithReport is like NormalizeEnv but also returns a report of
// changes made so callers can surface them to the user.
func NormalizeEnvWithReport(env map[string]string, opts NormalizeOptions) (map[string]string, NormalizeReport) {
	out := make(map[string]string, len(env))
	var report NormalizeReport

	// stable iteration order for deterministic report output
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]
		newKey := k

		if opts.TrimSpace {
			newV := strings.TrimSpace(v)
			if newV != v {
				report.Trimmed = append(report.Trimmed, k)
			}
			v = newV
			newKey = strings.TrimSpace(k)
		}
		if opts.LowercaseKey {
			lk := strings.ToLower(newKey)
			if lk != newKey {
				report.Renamed = append(report.Renamed, k)
			}
			newKey = lk
		}
		if opts.LowercaseVal {
			v = strings.ToLower(v)
		}
		if opts.CollapseEmpty && v == "" {
			report.Dropped = append(report.Dropped, k)
			continue
		}
		out[newKey] = v
	}
	return out, report
}
