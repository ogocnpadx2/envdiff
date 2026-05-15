package diff

import (
	"fmt"
	"sort"
	"strings"
)

// PivotEntry represents a single key's values across multiple environments.
type PivotEntry struct {
	Key    string
	Values map[string]string // env label -> value
}

// PivotTable holds a cross-environment view of all keys.
type PivotTable struct {
	Envs    []string      // ordered environment labels
	Entries []PivotEntry  // one per unique key, sorted
}

// BuildPivot constructs a PivotTable from a map of env label -> parsed key/value map.
func BuildPivot(envs map[string]map[string]string) PivotTable {
	// collect ordered env labels
	labels := make([]string, 0, len(envs))
	for label := range envs {
		labels = append(labels, label)
	}
	sort.Strings(labels)

	// collect all unique keys
	keySet := map[string]struct{}{}
	for _, kv := range envs {
		for k := range kv {
			keySet[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]PivotEntry, 0, len(keys))
	for _, k := range keys {
		values := make(map[string]string, len(labels))
		for _, label := range labels {
			if v, ok := envs[label][k]; ok {
				values[label] = v
			} else {
				values[label] = ""
			}
		}
		entries = append(entries, PivotEntry{Key: k, Values: values})
	}

	return PivotTable{Envs: labels, Entries: entries}
}

// FormatPivotText renders the PivotTable as an aligned text table.
func FormatPivotText(pt PivotTable) string {
	if len(pt.Entries) == 0 {
		return "(no keys)\n"
	}

	colWidth := 20
	var sb strings.Builder

	// header
	sb.WriteString(fmt.Sprintf("%-30s", "KEY"))
	for _, env := range pt.Envs {
		sb.WriteString(fmt.Sprintf(" %-*s", colWidth, truncate(env, colWidth)))
	}
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("-", 30+len(pt.Envs)*(colWidth+1)) + "\n")

	for _, entry := range pt.Entries {
		sb.WriteString(fmt.Sprintf("%-30s", truncate(entry.Key, 30)))
		for _, env := range pt.Envs {
			v := entry.Values[env]
			if v == "" {
				v = "<missing>"
			}
			sb.WriteString(fmt.Sprintf(" %-*s", colWidth, truncate(v, colWidth)))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
