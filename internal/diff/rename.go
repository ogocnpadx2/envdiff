package diff

import "sort"

// RenameMap maps old key names to new key names.
type RenameMap map[string]string

// RenameEntry represents a single detected or applied rename.
type RenameEntry struct {
	OldKey string
	NewKey string
}

// RenameResult holds the outcome of applying a rename map to a diff result.
type RenameResult struct {
	Applied  []RenameEntry
	Skipped  []RenameEntry // new key not present in either env
	Conflicts []RenameEntry // old key not found as missing
}

// ApplyRenames takes a diff Result and a RenameMap and attempts to reconcile
// keys that were renamed between environments. Renamed keys are removed from
// MissingInRight / MissingInLeft and recorded as applied renames.
func ApplyRenames(r Result, renames RenameMap) (Result, RenameResult) {
	missingRight := toSet(r.MissingInRight)
	missingLeft := toSet(r.MissingInLeft)

	var rr RenameResult

	for oldKey, newKey := range renames {
		entry := RenameEntry{OldKey: oldKey, NewKey: newKey}
		// old key missing in right, new key missing in left => rename detected
		if missingRight[oldKey] && missingLeft[newKey] {
			delete(missingRight, oldKey)
			delete(missingLeft, newKey)
			rr.Applied = append(rr.Applied, entry)
		} else if !missingRight[oldKey] && !missingLeft[newKey] {
			rr.Conflicts = append(rr.Conflicts, entry)
		} else {
			rr.Skipped = append(rr.Skipped, entry)
		}
	}

	r.MissingInRight = setToSlice(missingRight)
	r.MissingInLeft = setToSlice(missingLeft)
	sortStrings(r.MissingInRight)
	sortStrings(r.MissingInLeft)

	sortRenameEntries(rr.Applied)
	sortRenameEntries(rr.Skipped)
	sortRenameEntries(rr.Conflicts)

	return r, rr
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}

func setToSlice(s map[string]bool) []string {
	out := make([]string, 0, len(s))
	for k := range s {
		out = append(out, k)
	}
	return out
}

func sortRenameEntries(entries []RenameEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].OldKey < entries[j].OldKey
	})
}
