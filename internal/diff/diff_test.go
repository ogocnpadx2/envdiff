package diff_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/diff"
)

func TestCompare_Clean(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1", "B": "2"}

	result := diff.Compare(left, right)
	if !result.IsClean() {
		t.Errorf("expected clean result, got %+v", result)
	}
}

func TestCompare_MissingInRight(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1"}

	result := diff.Compare(left, right)
	if len(result.MissingInRight) != 1 || result.MissingInRight[0] != "B" {
		t.Errorf("expected B missing in right, got %v", result.MissingInRight)
	}
}

func TestCompare_MissingInLeft(t *testing.T) {
	left := map[string]string{"A": "1"}
	right := map[string]string{"A": "1", "C": "3"}

	result := diff.Compare(left, right)
	if len(result.MissingInLeft) != 1 || result.MissingInLeft[0] != "C" {
		t.Errorf("expected C missing in left, got %v", result.MissingInLeft)
	}
}

func TestCompare_Mismatched(t *testing.T) {
	left := map[string]string{"A": "foo", "B": "bar"}
	right := map[string]string{"A": "foo", "B": "baz"}

	result := diff.Compare(left, right)
	if len(result.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(result.Mismatched))
	}
	m := result.Mismatched[0]
	if m.Key != "B" || m.LeftValue != "bar" || m.RightValue != "baz" {
		t.Errorf("unexpected mismatch entry: %+v", m)
	}
}

func TestCompare_SortedOutput(t *testing.T) {
	left := map[string]string{"Z": "1", "A": "1", "M": "1"}
	right := map[string]string{}

	result := diff.Compare(left, right)
	expected := []string{"A", "M", "Z"}
	for i, k := range result.MissingInRight {
		if k != expected[i] {
			t.Errorf("expected sorted key %s at index %d, got %s", expected[i], i, k)
		}
	}
}

func TestCompare_BothEmpty(t *testing.T) {
	result := diff.Compare(map[string]string{}, map[string]string{})
	if !result.IsClean() {
		t.Error("expected clean result for two empty maps")
	}
}
