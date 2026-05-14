package diff

import (
	"testing"
)

func TestFindDuplicates_NoDuplicates(t *testing.T) {
	lines := []string{"FOO=bar", "BAZ=qux"}
	res := FindDuplicates(lines, DefaultDedupeOptions())
	if len(res.Duplicates) != 0 {
		t.Fatalf("expected no duplicates, got %v", res.Duplicates)
	}
	if res.Clean["FOO"] != "bar" || res.Clean["BAZ"] != "qux" {
		t.Fatalf("unexpected clean map: %v", res.Clean)
	}
}

func TestFindDuplicates_SingleDuplicate(t *testing.T) {
	lines := []string{"FOO=first", "BAR=ok", "FOO=second"}
	res := FindDuplicates(lines, DefaultDedupeOptions())
	if len(res.Duplicates) != 1 {
		t.Fatalf("expected 1 duplicate, got %d", len(res.Duplicates))
	}
	d := res.Duplicates[0]
	if d.Key != "FOO" {
		t.Errorf("expected key FOO, got %s", d.Key)
	}
	if len(d.Values) != 2 {
		t.Errorf("expected 2 values, got %v", d.Values)
	}
	// clean map should hold last-seen value
	if res.Clean["FOO"] != "second" {
		t.Errorf("expected clean value 'second', got %s", res.Clean["FOO"])
	}
}

func TestFindDuplicates_MultipleDuplicates(t *testing.T) {
	lines := []string{"A=1", "B=2", "A=3", "B=4", "B=5"}
	res := FindDuplicates(lines, DefaultDedupeOptions())
	if len(res.Duplicates) != 2 {
		t.Fatalf("expected 2 duplicates, got %d", len(res.Duplicates))
	}
	// sorted by key: A then B
	if res.Duplicates[0].Key != "A" || res.Duplicates[1].Key != "B" {
		t.Errorf("unexpected order: %v", res.Duplicates)
	}
	if res.Clean["B"] != "5" {
		t.Errorf("expected last value '5' for B, got %s", res.Clean["B"])
	}
}

func TestFindDuplicates_CaseInsensitive(t *testing.T) {
	lines := []string{"FOO=upper", "foo=lower"}
	opts := DedupeOptions{CaseSensitive: false}
	res := FindDuplicates(lines, opts)
	if len(res.Duplicates) != 1 {
		t.Fatalf("expected 1 duplicate (case-insensitive), got %d", len(res.Duplicates))
	}
	if res.Duplicates[0].Key != "foo" {
		t.Errorf("expected normalised key 'foo', got %s", res.Duplicates[0].Key)
	}
}

func TestFindDuplicates_SkipsInvalidLines(t *testing.T) {
	lines := []string{"# comment", "", "VALID=yes", "no-equals"}
	res := FindDuplicates(lines, DefaultDedupeOptions())
	if len(res.Duplicates) != 0 {
		t.Fatalf("expected no duplicates, got %v", res.Duplicates)
	}
	if _, ok := res.Clean["VALID"]; !ok {
		t.Error("expected VALID key in clean map")
	}
}

func TestFindDuplicates_EmptyInput(t *testing.T) {
	res := FindDuplicates([]string{}, DefaultDedupeOptions())
	if len(res.Duplicates) != 0 {
		t.Error("expected no duplicates for empty input")
	}
	if len(res.Clean) != 0 {
		t.Error("expected empty clean map for empty input")
	}
}
