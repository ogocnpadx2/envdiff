package diff

import (
	"fmt"
	"sort"
	"strings"
)

// AliasMap maps canonical key names to one or more alternate names.
type AliasMap map[string][]string

// AliasResult holds the outcome of resolving aliases against an env map.
type AliasResult struct {
	// Resolved maps canonical key -> value found via alias
	Resolved map[string]string
	// Unresolved lists canonical keys for which no alias matched
	Unresolved []string
	// UsedAlias maps canonical key -> alias that was matched
	UsedAlias map[string]string
}

// ParseAliasMap parses a slice of "canonical=alias1,alias2" strings into an AliasMap.
func ParseAliasMap(entries []string) (AliasMap, error) {
	am := make(AliasMap)
	for _, entry := range entries {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid alias entry %q: expected canonical=alias1,alias2", entry)
		}
		canonical := strings.TrimSpace(parts[0])
		aliases := strings.Split(parts[1], ",")
		for i, a := range aliases {
			aliases[i] = strings.TrimSpace(a)
		}
		am[canonical] = aliases
	}
	return am, nil
}

// ResolveAliases looks up canonical keys in env, falling back to aliases when
// the canonical key is absent. Keys already present in env are not overridden.
func ResolveAliases(env map[string]string, am AliasMap) AliasResult {
	resolved := make(map[string]string)
	usedAlias := make(map[string]string)
	var unresolved []string

	canonicals := make([]string, 0, len(am))
	for k := range am {
		canonicals = append(canonicals, k)
	}
	sort.Strings(canonicals)

	for _, canonical := range canonicals {
		if v, ok := env[canonical]; ok {
			resolved[canonical] = v
			continue
		}
		found := false
		for _, alias := range am[canonical] {
			if v, ok := env[alias]; ok {
				resolved[canonical] = v
				usedAlias[canonical] = alias
				found = true
				break
			}
		}
		if !found {
			unresolved = append(unresolved, canonical)
		}
	}

	return AliasResult{
		Resolved:   resolved,
		Unresolved: unresolved,
		UsedAlias:  usedAlias,
	}
}
