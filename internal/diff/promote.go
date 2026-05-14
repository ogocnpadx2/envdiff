package diff

import (
	"fmt"
	"sort"
)

// PromoteStrategy controls how values are selected when promoting.
type PromoteStrategy int

const (
	PromoteOnlyMissing PromoteStrategy = iota // only fill in keys absent in target
	PromoteOverwrite                          // overwrite all keys from source into target
)

// PromoteOptions configures a promotion run.
type PromoteOptions struct {
	Strategy PromoteStrategy
	DryRun   bool
}

// PromoteResult captures what happened during promotion.
type PromoteResult struct {
	Added     []string          // keys added to target
	Overwritten []string        // keys overwritten in target
	Skipped   []string          // keys skipped (already present, OnlyMissing mode)
	Merged    map[string]string // final merged env map
}

// ParsePromoteStrategy converts a string flag to a PromoteStrategy.
func ParsePromoteStrategy(s string) (PromoteStrategy, error) {
	switch s {
	case "missing", "":
		return PromoteOnlyMissing, nil
	case "overwrite":
		return PromoteOverwrite, nil
	default:
		return 0, fmt.Errorf("unknown promote strategy %q: want \"missing\" or \"overwrite\"", s)
	}
}

// Promote copies keys from source into target according to opts.
// It never mutates the input maps.
func Promote(source, target map[string]string, opts PromoteOptions) PromoteResult {
	merged := make(map[string]string, len(target))
	for k, v := range target {
		merged[k] = v
	}

	var added, overwritten, skipped []string

	for k, v := range source {
		existing, exists := target[k]
		switch {
		case !exists:
			added = append(added, k)
			if !opts.DryRun {
				merged[k] = v
			}
		case opts.Strategy == PromoteOverwrite && existing != v:
			overwritten = append(overwritten, k)
			if !opts.DryRun {
				merged[k] = v
			}
		default:
			skipped = append(skipped, k)
		}
	}

	sort.Strings(added)
	sort.Strings(overwritten)
	sort.Strings(skipped)

	return PromoteResult{
		Added:       added,
		Overwritten: overwritten,
		Skipped:     skipped,
		Merged:      merged,
	}
}
