package diff

import "sort"

// CascadeStrategy controls how values flow between environments.
type CascadeStrategy int

const (
	CascadeOverwrite CascadeStrategy = iota // later envs overwrite earlier
	CascadePreserve                         // earlier envs take precedence
	CascadeMissing                          // only fill in missing keys
)

// CascadeEntry holds a key and the resolved value along with its origin.
type CascadeEntry struct {
	Key    string
	Value  string
	Origin string // label/path of the env that provided this value
}

// CascadeResult is the output of a cascade operation.
type CascadeResult struct {
	Resolved []CascadeEntry
	Skipped  []CascadeEntry // entries that were shadowed
}

// ParseCascadeStrategy parses a strategy name string.
func ParseCascadeStrategy(s string) (CascadeStrategy, error) {
	switch s {
	case "overwrite", "":
		return CascadeOverwrite, nil
	case "preserve":
		return CascadePreserve, nil
	case "missing":
		return CascadeMissing, nil
	}
	return 0, fmt.Errorf("unknown cascade strategy %q: want overwrite|preserve|missing", s)
}

// Cascade merges multiple labeled env maps according to the given strategy.
// envs is ordered from lowest to highest priority (index 0 = base).
func Cascade(envs []map[string]string, labels []string, strategy CascadeStrategy) CascadeResult {
	resolved := map[string]CascadeEntry{}
	var skipped []CascadeEntry

	for i, env := range envs {
		label := ""
		if i < len(labels) {
			label = labels[i]
		}
		for k, v := range env {
			existing, exists := resolved[k]
			switch strategy {
			case CascadeOverwrite:
				if exists {
					skipped = append(skipped, existing)
				}
				resolved[k] = CascadeEntry{Key: k, Value: v, Origin: label}
			case CascadePreserve:
				if !exists {
					resolved[k] = CascadeEntry{Key: k, Value: v, Origin: label}
				} else {
					skipped = append(skipped, CascadeEntry{Key: k, Value: v, Origin: label})
				}
			case CascadeMissing:
				if !exists {
					resolved[k] = CascadeEntry{Key: k, Value: v, Origin: label}
				} else {
					skipped = append(skipped, CascadeEntry{Key: k, Value: v, Origin: label})
				}
			}
		}
	}

	out := make([]CascadeEntry, 0, len(resolved))
	for _, e := range resolved {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	sort.Slice(skipped, func(i, j int) bool { return skipped[i].Key < skipped[j].Key })
	return CascadeResult{Resolved: out, Skipped: skipped}
}
