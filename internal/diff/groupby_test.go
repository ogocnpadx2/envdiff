package diff

import (
	"testing"
)

func baseGroupResult() Result {
	return Result{
		MissingInRight: []string{"DB_HOST", "DB_PORT", "APP_NAME"},
		MissingInLeft:  []string{"AWS_KEY"},
		Mismatched: []Mismatch{
			{Key: "DB_PASS", LeftVal: "secret", RightVal: "other"},
			{Key: "TIMEOUT", LeftVal: "30", RightVal: "60"},
		},
	}
}

func TestGroupResult_PrefixBuckets(t *testing.T) {
	groups := GroupResult(baseGroupResult())

	prefixMap := make(map[string][]string)
	for _, g := range groups {
		prefixMap[g.Prefix] = g.Keys
	}

	if _, ok := prefixMap["DB_"]; !ok {
		t.Fatal("expected DB_ prefix group")
	}
	if _, ok := prefixMap["AWS_"]; !ok {
		t.Fatal("expected AWS_ prefix group")
	}
	if _, ok := prefixMap["APP_"]; !ok {
		t.Fatal("expected APP_ prefix group")
	}
	if _, ok := prefixMap[""]; !ok {
		t.Fatal("expected empty prefix group for keys without underscore")
	}
}

func TestGroupResult_DBKeys(t *testing.T) {
	groups := GroupResult(baseGroupResult())
	for _, g := range groups {
		if g.Prefix != "DB_" {
			continue
		}
		if len(g.Keys) != 3 {
			t.Fatalf("expected 3 DB_ keys, got %d: %v", len(g.Keys), g.Keys)
		}
		return
	}
	t.Fatal("DB_ group not found")
}

func TestGroupResult_EmptyResult(t *testing.T) {
	groups := GroupResult(Result{})
	if len(groups) != 0 {
		t.Fatalf("expected no groups for empty result, got %d", len(groups))
	}
}

func TestGroupResult_NoDuplicates(t *testing.T) {
	// A key appearing in both MissingInRight and Mismatched (edge case).
	r := Result{
		MissingInRight: []string{"DB_HOST", "DB_HOST"},
		Mismatched:     []Mismatch{{Key: "DB_HOST", LeftVal: "a", RightVal: "b"}},
	}
	groups := GroupResult(r)
	for _, g := range groups {
		if g.Prefix == "DB_" {
			if len(g.Keys) != 1 {
				t.Fatalf("expected 1 unique key, got %d: %v", len(g.Keys), g.Keys)
			}
			return
		}
	}
	t.Fatal("DB_ group not found")
}

func TestExtractPrefix(t *testing.T) {
	cases := []struct {
		key    string
		want   string
	}{
		{"DB_HOST", "DB_"},
		{"AWS_SECRET_KEY", "AWS_"},
		{"TIMEOUT", ""},
		{"_LEADING", "_"},
	}
	for _, tc := range cases {
		got := extractPrefix(tc.key)
		if got != tc.want {
			t.Errorf("extractPrefix(%q) = %q, want %q", tc.key, got, tc.want)
		}
	}
}
