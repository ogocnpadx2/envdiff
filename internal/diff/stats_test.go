package diff

import (
	"bytes"
	"strings"
	"testing"
)

func baseStats() Result {
	return Result{
		MissingInLeft:  []string{"A"},
		MissingInRight: []string{"B", "C"},
		Mismatched: []MismatchedKey{
			{Key: "D", LeftVal: "x", RightVal: "y"},
		},
	}
}

func TestComputeStats_Clean(t *testing.T) {
	s := ComputeStats(Result{})
	if !s.IsClean() {
		t.Error("expected clean")
	}
	if s.Total != 0 {
		t.Errorf("expected total 0, got %d", s.Total)
	}
}

func TestComputeStats_WithDiffs(t *testing.T) {
	s := ComputeStats(baseStats())
	if s.MissingInLeft != 1 {
		t.Errorf("MissingInLeft: want 1, got %d", s.MissingInLeft)
	}
	if s.MissingInRight != 2 {
		t.Errorf("MissingInRight: want 2, got %d", s.MissingInRight)
	}
	if s.Mismatched != 1 {
		t.Errorf("Mismatched: want 1, got %d", s.Mismatched)
	}
	if s.Total != 4 {
		t.Errorf("Total: want 4, got %d", s.Total)
	}
}

func TestStats_IsClean(t *testing.T) {
	if !(Stats{}).IsClean() {
		t.Error("zero Stats should be clean")
	}
	if (Stats{Total: 1}).IsClean() {
		t.Error("non-zero Stats should not be clean")
	}
}

func TestStats_Summary_Clean(t *testing.T) {
	s := Stats{}
	if s.Summary() != "No differences found." {
		t.Errorf("unexpected summary: %s", s.Summary())
	}
}

func TestStats_Summary_WithDiffs(t *testing.T) {
	s := ComputeStats(baseStats())
	got := s.Summary()
	if !strings.Contains(got, "4 issue(s)") {
		t.Errorf("summary missing issue count: %s", got)
	}
}

func TestPrintStats_Clean(t *testing.T) {
	var buf bytes.Buffer
	PrintStats(&buf, Result{})
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestPrintStats_WithDiffs(t *testing.T) {
	var buf bytes.Buffer
	PrintStats(&buf, baseStats())
	out := buf.String()
	for _, want := range []string{"missing in left", "missing in right", "mismatched"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q: %s", want, out)
		}
	}
}
