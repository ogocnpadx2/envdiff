package diff

import (
	"testing"
)

func TestParseSortField_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  SortField
	}{
		{"key", SortByKey},
		{"value", SortByValue},
		{"length", SortByLength},
		{"KEY", SortByKey},
	}
	for _, tc := range cases {
		got, err := ParseSortField(tc.input)
		if err != nil {
			t.Errorf("ParseSortField(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseSortField(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestParseSortField_Invalid(t *testing.T) {
	_, err := ParseSortField("unknown")
	if err == nil {
		t.Error("expected error for unknown sort field")
	}
}

func TestParseSortOrder_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  SortOrder
	}{
		{"asc", SortAsc},
		{"desc", SortDesc},
		{"", SortAsc},
	}
	for _, tc := range cases {
		got, err := ParseSortOrder(tc.input)
		if err != nil {
			t.Errorf("ParseSortOrder(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseSortOrder(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestParseSortOrder_Invalid(t *testing.T) {
	_, err := ParseSortOrder("sideways")
	if err == nil {
		t.Error("expected error for unknown sort order")
	}
}

func TestSortEnv_ByKeyAsc(t *testing.T) {
	env := map[string]string{"ZEBRA": "z", "APPLE": "a", "MANGO": "m"}
	result := SortEnv(env, DefaultSortOptions())
	keys := []string{result[0].Key, result[1].Key, result[2].Key}
	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("position %d: got %q, want %q", i, keys[i], k)
		}
	}
}

func TestSortEnv_ByKeyDesc(t *testing.T) {
	env := map[string]string{"ZEBRA": "z", "APPLE": "a", "MANGO": "m"}
	opts := SortOptions{Field: SortByKey, Order: SortDesc}
	result := SortEnv(env, opts)
	if result[0].Key != "ZEBRA" {
		t.Errorf("expected ZEBRA first in desc, got %q", result[0].Key)
	}
}

func TestSortEnv_ByValue(t *testing.T) {
	env := map[string]string{"A": "charlie", "B": "alpha", "C": "bravo"}
	opts := SortOptions{Field: SortByValue, Order: SortAsc}
	result := SortEnv(env, opts)
	if result[0].Value != "alpha" {
		t.Errorf("expected alpha first, got %q", result[0].Value)
	}
}

func TestSortEnv_ByLength(t *testing.T) {
	env := map[string]string{"A": "hi", "B": "hello world", "C": "hey"}
	opts := SortOptions{Field: SortByLength, Order: SortAsc}
	result := SortEnv(env, opts)
	if result[0].Value != "hi" {
		t.Errorf("expected shortest value first, got %q", result[0].Value)
	}
}

func TestSortEnv_IgnoreCase(t *testing.T) {
	env := map[string]string{"banana": "1", "Apple": "2", "cherry": "3"}
	opts := SortOptions{Field: SortByKey, Order: SortAsc, IgnoreCase: true}
	result := SortEnv(env, opts)
	if result[0].Key != "Apple" {
		t.Errorf("expected Apple first (case-insensitive), got %q", result[0].Key)
	}
}

func TestSortEnv_Empty(t *testing.T) {
	result := SortEnv(map[string]string{}, DefaultSortOptions())
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}
