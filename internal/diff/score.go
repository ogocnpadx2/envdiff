package diff

import "fmt"

// Score represents a numeric health score for an env comparison result.
type Score struct {
	Total     int
	Penalty   int
	Value     float64 // 0.0 to 100.0
	Grade     string
}

// ScoreOptions controls how penalties are weighted.
type ScoreOptions struct {
	MissingPenalty    int
	MismatchedPenalty int
}

// DefaultScoreOptions returns sensible defaults.
func DefaultScoreOptions() ScoreOptions {
	return ScoreOptions{
		MissingPenalty:    10,
		MismatchedPenalty: 5,
	}
}

// ComputeScore calculates a 0–100 health score for a Result.
// Missing keys are penalised more heavily than mismatched values.
func ComputeScore(r Result, opts ScoreOptions) Score {
	total := len(r.OnlyInLeft) + len(r.OnlyInRight) + len(r.Mismatched) + len(r.Matching)
	if total == 0 {
		return Score{Total: 0, Penalty: 0, Value: 100.0, Grade: "A"}
	}

	penalty := 0
	penalty += len(r.OnlyInLeft) * opts.MissingPenalty
	penalty += len(r.OnlyInRight) * opts.MissingPenalty
	penalty += len(r.Mismatched) * opts.MismatchedPenalty

	// Cap penalty at total * max-possible-penalty so score stays >= 0
	maxPenalty := total * opts.MissingPenalty
	if penalty > maxPenalty {
		penalty = maxPenalty
	}

	raw := 100.0 * (1.0 - float64(penalty)/float64(maxPenalty))
	if raw < 0 {
		raw = 0
	}

	return Score{
		Total:   total,
		Penalty: penalty,
		Value:   raw,
		Grade:   scoreGrade(raw),
	}
}

func scoreGrade(v float64) string {
	switch {
	case v >= 90:
		return "A"
	case v >= 75:
		return "B"
	case v >= 60:
		return "C"
	case v >= 40:
		return "D"
	default:
		return "F"
	}
}

// Summary returns a one-line human-readable summary of the score.
func (s Score) Summary() string {
	return fmt.Sprintf("Score: %.1f/100 (Grade: %s, Penalty: %d, Keys: %d)",
		s.Value, s.Grade, s.Penalty, s.Total)
}
