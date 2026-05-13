package diff

import (
	"sort"
	"strings"
)

// SuggestionEntry holds a key that is missing in one env and the closest
// matching key found in the other env, based on edit distance.
type SuggestionEntry struct {
	MissingKey  string
	SuggestedKey string
	Distance    int
}

// SuggestResult contains rename/typo suggestions derived from a CompareResult.
type SuggestResult struct {
	Suggestions []SuggestionEntry
}

// Suggest analyses a CompareResult and proposes likely typo/rename matches
// for keys that appear only on one side, using Levenshtein distance.
func Suggest(result CompareResult, maxDistance int) SuggestResult {
	if maxDistance <= 0 {
		maxDistance = 3
	}

	var entries []SuggestionEntry

	candidates := append(append([]string{}, result.MissingInRight...), result.MissingInLeft...)

	for _, missing := range result.MissingInRight {
		best, dist := closestKey(missing, result.MissingInLeft)
		if dist <= maxDistance && best != "" {
			entries = append(entries, SuggestionEntry{
				MissingKey:   missing,
				SuggestedKey: best,
				Distance:     dist,
			})
		}
	}
	_ = candidates

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Distance != entries[j].Distance {
			return entries[i].Distance < entries[j].Distance
		}
		return entries[i].MissingKey < entries[j].MissingKey
	})

	return SuggestResult{Suggestions: entries}
}

func closestKey(target string, candidates []string) (string, int) {
	best := ""
	bestDist := int(^uint(0) >> 1)
	for _, c := range candidates {
		d := levenshtein(strings.ToLower(target), strings.ToLower(c))
		if d < bestDist {
			bestDist = d
			best = c
		}
	}
	return best, bestDist
}

func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}
	row := make([]int, lb+1)
	for j := 0; j <= lb; j++ {
		row[j] = j
	}
	for i := 1; i <= la; i++ {
		prev := row[0]
		row[0] = i
		for j := 1; j <= lb; j++ {
			old := row[j]
			if ra[i-1] == rb[j-1] {
				row[j] = prev
			} else {
				row[j] = 1 + min3(prev, row[j], row[j-1])
			}
			prev = old
		}
	}
	return row[lb]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
