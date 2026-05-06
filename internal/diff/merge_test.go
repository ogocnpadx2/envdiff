package diff

import (
	"testing"
)

func TestParseMergeStrategy_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected MergeStrategy
	}{
		{"left", MergePreferLeft},
		{"right", MergePreferRight},
		{"union", MergeUnionAll},
	}
	for _, tc := range cases {
		s, err := ParseMergeStrategy(tc.input)
		if err != nil {
			t.Errorf("unexpected error for %q: %v", tc.input, err)
		}
		if s != tc.expected {
			t.Errorf("expected %d, got %d", tc.expected, s)
		}
	}
}

func TestParseMergeStrategy_Invalid(t *testing.T) {
	_, err := ParseMergeStrategy("bogus")
	if err == nil {
		t.Error("expected error for invalid strategy, got nil")
	}
}

func TestMerge_PreferLeft(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"B": "99", "C": "3"}

	result := Merge(left, right, MergePreferLeft)

	if result["A"] != "1" {
		t.Errorf("expected A=1, got %s", result["A"])
	}
	if result["B"] != "2" {
		t.Errorf("expected B=2 (left wins), got %s", result["B"])
	}
	if result["C"] != "3" {
		t.Errorf("expected C=3, got %s", result["C"])
	}
}

func TestMerge_PreferRight(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"B": "99", "C": "3"}

	result := Merge(left, right, MergePreferRight)

	if result["B"] != "99" {
		t.Errorf("expected B=99 (right wins), got %s", result["B"])
	}
	if result["A"] != "1" {
		t.Errorf("expected A=1, got %s", result["A"])
	}
}

func TestMerge_UnionAll_ContainsAllKeys(t *testing.T) {
	left := map[string]string{"X": "left"}
	right := map[string]string{"Y": "right", "X": "right"}

	result := Merge(left, right, MergeUnionAll)

	if _, ok := result["X"]; !ok {
		t.Error("expected key X in union result")
	}
	if _, ok := result["Y"]; !ok {
		t.Error("expected key Y in union result")
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestMergedKeys_Sorted(t *testing.T) {
	m := map[string]string{"Z": "1", "A": "2", "M": "3"}
	keys := MergedKeys(m)
	expected := []string{"A", "M", "Z"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: expected %s, got %s", i, expected[i], k)
		}
	}
}
