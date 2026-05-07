package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline represents a saved reference state of a diff result.
type Baseline struct {
	CreatedAt   time.Time         `json:"created_at"`
	LeftFile    string            `json:"left_file"`
	RightFile   string            `json:"right_file"`
	MissingKeys map[string]string `json:"missing_keys"` // key -> which side is missing ("left" or "right")
	Mismatched  map[string][2]string `json:"mismatched"` // key -> [leftVal, rightVal]
}

// SaveBaseline writes the current diff result as a baseline JSON file.
func SaveBaseline(path, leftFile, rightFile string, result Result) error {
	b := Baseline{
		CreatedAt:   time.Now().UTC(),
		LeftFile:    leftFile,
		RightFile:   rightFile,
		MissingKeys: make(map[string]string),
		Mismatched:  make(map[string][2]string),
	}
	for _, k := range result.MissingInLeft {
		b.MissingKeys[k] = "left"
	}
	for _, k := range result.MissingInRight {
		b.MissingKeys[k] = "right"
	}
	for _, m := range result.Mismatched {
		b.Mismatched[m.Key] = [2]string{m.LeftVal, m.RightVal}
	}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadBaseline reads a baseline from disk.
func LoadBaseline(path string) (Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Baseline{}, fmt.Errorf("baseline: read %s: %w", path, err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return Baseline{}, fmt.Errorf("baseline: parse: %w", err)
	}
	return b, nil
}

// DiffAgainstBaseline compares a current Result against a saved Baseline.
// It returns new issues not present in the baseline and resolved issues that
// were in the baseline but are no longer present.
func DiffAgainstBaseline(b Baseline, current Result) (newIssues Result, resolved Result) {
	currentMissing := map[string]string{}
	for _, k := range current.MissingInLeft {
		currentMissing[k] = "left"
	}
	for _, k := range current.MissingInRight {
		currentMissing[k] = "right"
	}

	// New missing keys not in baseline
	for k, side := range currentMissing {
		if _, known := b.MissingKeys[k]; !known {
			if side == "left" {
				newIssues.MissingInLeft = append(newIssues.MissingInLeft, k)
			} else {
				newIssues.MissingInRight = append(newIssues.MissingInRight, k)
			}
		}
	}

	// Resolved missing keys present in baseline but gone now
	for k, side := range b.MissingKeys {
		if _, still := currentMissing[k]; !still {
			if side == "left" {
				resolved.MissingInLeft = append(resolved.MissingInLeft, k)
			} else {
				resolved.MissingInRight = append(resolved.MissingInRight, k)
			}
		}
	}

	currentMM := map[string]MismatchedKey{}
	for _, m := range current.Mismatched {
		currentMM[m.Key] = m
	}

	// New mismatches not in baseline
	for _, m := range current.Mismatched {
		if _, known := b.Mismatched[m.Key]; !known {
			newIssues.Mismatched = append(newIssues.Mismatched, m)
		}
	}

	// Resolved mismatches
	for k, vals := range b.Mismatched {
		if cur, still := currentMM[k]; !still || cur.LeftVal != vals[0] || cur.RightVal != vals[1] {
			if !still {
				resolved.Mismatched = append(resolved.Mismatched, MismatchedKey{Key: k, LeftVal: vals[0], RightVal: vals[1]})
			}
		}
	}

	return newIssues, resolved
}
