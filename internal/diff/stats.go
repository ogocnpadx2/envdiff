package diff

import "fmt"

// Stats holds summary counts for a diff result.
type Stats struct {
	MissingInLeft  int
	MissingInRight int
	Mismatched     int
	Total          int
}

// ComputeStats derives Stats from a Result.
func ComputeStats(r Result) Stats {
	return Stats{
		MissingInLeft:  len(r.MissingInLeft),
		MissingInRight: len(r.MissingInRight),
		Mismatched:     len(r.Mismatched),
		Total:          len(r.MissingInLeft) + len(r.MissingInRight) + len(r.Mismatched),
	}
}

// IsClean returns true when there are no differences.
func (s Stats) IsClean() bool {
	return s.Total == 0
}

// Summary returns a human-readable one-line summary.
func (s Stats) Summary() string {
	if s.IsClean() {
		return "No differences found."
	}
	return fmt.Sprintf(
		"%d issue(s): %d missing in left, %d missing in right, %d mismatched",
		s.Total, s.MissingInLeft, s.MissingInRight, s.Mismatched,
	)
}

// PrintStats writes a stats summary to w.
func PrintStats(w interface{ Write([]byte) (int, error) }, r Result) {
	s := ComputeStats(r)
	fmt.Fprintln(w, s.Summary())
	if !s.IsClean() {
		fmt.Fprintf(w, "  missing in left  : %d\n", s.MissingInLeft)
		fmt.Fprintf(w, "  missing in right : %d\n", s.MissingInRight)
		fmt.Fprintf(w, "  mismatched       : %d\n", s.Mismatched)
	}
}
