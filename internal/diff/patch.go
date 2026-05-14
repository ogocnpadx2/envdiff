package diff

import (
	"fmt"
	"sort"
	"strings"
)

// PatchOp represents a single patch operation type.
type PatchOp string

const (
	PatchAdd    PatchOp = "add"
	PatchRemove PatchOp = "remove"
	PatchChange PatchOp = "change"
)

// PatchEntry describes a single change to apply to an env map.
type PatchEntry struct {
	Op       PatchOp
	Key      string
	OldValue string // populated for change/remove
	NewValue string // populated for add/change
}

// PatchResult holds the outcome of applying a patch.
type PatchResult struct {
	Applied  []PatchEntry
	Skipped  []PatchEntry
	Conflicts []PatchEntry
}

// BuildPatch computes the patch (set of operations) needed to transform
// src into dst.
func BuildPatch(src, dst map[string]string) []PatchEntry {
	var entries []PatchEntry

	for k, dv := range dst {
		if sv, ok := src[k]; !ok {
			entries = append(entries, PatchEntry{Op: PatchAdd, Key: k, NewValue: dv})
		} else if sv != dv {
			entries = append(entries, PatchEntry{Op: PatchChange, Key: k, OldValue: sv, NewValue: dv})
		}
	}

	for k, sv := range src {
		if _, ok := dst[k]; !ok {
			entries = append(entries, PatchEntry{Op: PatchRemove, Key: k, OldValue: sv})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Key != entries[j].Key {
			return entries[i].Key < entries[j].Key
		}
		return entries[i].Op < entries[j].Op
	})
	return entries
}

// ApplyPatch applies patch entries to target, returning a new map and a PatchResult.
// If dryRun is true, no modifications are made but the result is still populated.
func ApplyPatch(target map[string]string, entries []PatchEntry, dryRun bool) (map[string]string, PatchResult) {
	out := make(map[string]string, len(target))
	for k, v := range target {
		out[k] = v
	}

	var result PatchResult
	for _, e := range entries {
		switch e.Op {
		case PatchAdd:
			if _, exists := out[e.Key]; exists {
				result.Conflicts = append(result.Conflicts, e)
				continue
			}
			if !dryRun {
				out[e.Key] = e.NewValue
			}
			result.Applied = append(result.Applied, e)
		case PatchRemove:
			if _, exists := out[e.Key]; !exists {
				result.Skipped = append(result.Skipped, e)
				continue
			}
			if !dryRun {
				delete(out, e.Key)
			}
			result.Applied = append(result.Applied, e)
		case PatchChange:
			current, exists := out[e.Key]
			if !exists {
				result.Skipped = append(result.Skipped, e)
				continue
			}
			if current != e.OldValue {
				result.Conflicts = append(result.Conflicts, e)
				continue
			}
			if !dryRun {
				out[e.Key] = e.NewValue
			}
			result.Applied = append(result.Applied, e)
		}
	}
	return out, result
}

// FormatPatch renders a patch as a human-readable unified-diff-style string.
func FormatPatch(entries []PatchEntry) string {
	var sb strings.Builder
	for _, e := range entries {
		switch e.Op {
		case PatchAdd:
			fmt.Fprintf(&sb, "+ %s=%s\n", e.Key, e.NewValue)
		case PatchRemove:
			fmt.Fprintf(&sb, "- %s=%s\n", e.Key, e.OldValue)
		case PatchChange:
			fmt.Fprintf(&sb, "~ %s: %s -> %s\n", e.Key, e.OldValue, e.NewValue)
		}
	}
	return sb.String()
}
