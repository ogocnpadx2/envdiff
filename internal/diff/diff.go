package diff

// Result holds the comparison outcome between two env files.
type Result struct {
	// MissingInRight contains keys present in left but absent in right.
	MissingInRight []string
	// MissingInLeft contains keys present in right but absent in left.
	MissingInLeft []string
	// Mismatched contains keys present in both files but with different values.
	Mismatched []MismatchedKey
}

// MismatchedKey represents a key whose value differs between two env files.
type MismatchedKey struct {
	Key        string
	LeftValue  string
	RightValue string
}

// IsClean returns true when there are no differences between the two env maps.
func (r Result) IsClean() bool {
	return len(r.MissingInRight) == 0 &&
		len(r.MissingInLeft) == 0 &&
		len(r.Mismatched) == 0
}

// Compare takes two parsed env maps and returns a Result describing their differences.
func Compare(left, right map[string]string) Result {
	var result Result

	for k, lv := range left {
		rv, ok := right[k]
		if !ok {
			result.MissingInRight = append(result.MissingInRight, k)
			continue
		}
		if lv != rv {
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

	sortStrings(result.MissingInRight)
	sortStrings(result.MissingInLeft)
	sortMismatched(result.Mismatched)

	return result
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}

func sortMismatched(m []MismatchedKey) {
	for i := 1; i < len(m); i++ {
		for j := i; j > 0 && m[j].Key < m[j-1].Key; j-- {
			m[j], m[j-1] = m[j-1], m[j]
		}
	}
}
