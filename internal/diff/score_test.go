package diff

import (
	"strings"
	"testing"
)

func cleanScoreResult() Result {
	return Result{
		Matching:   []string{"A", "B", "C", "D", "E"},
		OnlyInLeft: []string{},
		OnlyInRight: []string{},
		Mismatched: []MismatchedKey{},
	}
}

func TestComputeScore_PerfectClean(t *testing.T) {
	r := cleanScoreResult()
	s := ComputeScore(r, DefaultScoreOptions())
	if s.Value != 100.0 {
		t.Errorf("expected 100.0, got %.1f", s.Value)
	}
	if s.Grade != "A" {
		t.Errorf("expected grade A, got %s", s.Grade)
	}
}

func TestComputeScore_EmptyResult(t *testing.T) {
	s := ComputeScore(Result{}, DefaultScoreOptions())
	if s.Value != 100.0 {
		t.Errorf("expected 100.0 for empty result, got %.1f", s.Value)
	}
	if s.Total != 0 {
		t.Errorf("expected total 0, got %d", s.Total)
	}
}

func TestComputeScore_AllMissing(t *testing.T) {
	r := Result{
		OnlyInLeft:  []string{"X", "Y"},
		OnlyInRight: []string{"Z"},
		Matching:    []string{},
		Mismatched:  []MismatchedKey{},
	}
	s := ComputeScore(r, DefaultScoreOptions())
	if s.Value != 0.0 {
		t.Errorf("expected 0.0 for all-missing, got %.1f", s.Value)
	}
	if s.Grade != "F" {
		t.Errorf("expected grade F, got %s", s.Grade)
	}
}

func TestComputeScore_PartialMismatch(t *testing.T) {
	r := Result{
		Matching:    []string{"A", "B", "C"},
		OnlyInLeft:  []string{},
		OnlyInRight: []string{},
		Mismatched:  []MismatchedKey{{Key: "D"}, {Key: "E"}},
	}
	s := ComputeScore(r, DefaultScoreOptions())
	if s.Value >= 100.0 {
		t.Errorf("expected score < 100 with mismatches, got %.1f", s.Value)
	}
	if s.Value <= 0.0 {
		t.Errorf("expected score > 0 with some matching, got %.1f", s.Value)
	}
}

func TestComputeScore_CustomPenalty(t *testing.T) {
	r := Result{
		Matching:    []string{"A"},
		OnlyInLeft:  []string{"B"},
		OnlyInRight: []string{},
		Mismatched:  []MismatchedKey{},
	}
	opts := ScoreOptions{MissingPenalty: 20, MismatchedPenalty: 2}
	s := ComputeScore(r, opts)
	if s.Value >= 100.0 {
		t.Errorf("expected penalty applied, got %.1f", s.Value)
	}
}

func TestScore_Summary(t *testing.T) {
	r := cleanScoreResult()
	s := ComputeScore(r, DefaultScoreOptions())
	sum := s.Summary()
	if !strings.Contains(sum, "100.0") {
		t.Errorf("expected summary to contain score, got: %s", sum)
	}
	if !strings.Contains(sum, "Grade: A") {
		t.Errorf("expected summary to contain grade, got: %s", sum)
	}
}
