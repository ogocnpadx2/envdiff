package diff

import "sort"

// MismatchedKey holds a key whose value differs between the two env files.
type MismatchedKey struct {
	Key        string
	LeftValue  string
	RightValue string
}

// Result holds the full comparison output.
type Result struct {
	MissingInLeft  []string
	MissingInRight []string
	Mismatched     []MismatchedKey
}

// IsClean returns true when there are no differences.
func (r Result) IsClean() bool {
	return len(r.MissingInLeft) == 0 &&
		len(r.MissingInRight) == 0 &&
		len(r.Mismatched) == 0
}

// Compare compares two maps of env variables and returns a Result.
func Compare(left, right map[string]string) Result {
	var result Result

	for k, lv := range left {
		rv, ok := right[k]
		if !ok {
			result.MissingInRight = append(result.MissingInRight, k)
		} else if lv != rv {
			result.Mismatched = append(result.Mismatched, MismatchedKey{
				Key:        k,
				LeftValue:  lv,
				RightValue: rv,
			})
		}
	}

	for k := range right {
		if _, ok := left[k]; !ok {
			result.MissingInLeft = append(result.MissingInLeft, k)
		}
	}

	sortStrings(result.MissingInLeft)
	sortStrings(result.MissingInRight)
	sortMismatched(result.Mismatched)

	return result
}

func sortStrings(s []string) {
	sort.Strings(s)
}

func sortMismatched(m []MismatchedKey) {
	sort.Slice(m, func(i, j int) bool {
		return m[i].Key < m[j].Key
	})
}
