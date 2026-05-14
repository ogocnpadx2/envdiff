package diff

import (
	"testing"
)

func TestParseCascadeStrategy_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  CascadeStrategy
	}{
		{"overwrite", CascadeOverwrite},
		{"", CascadeOverwrite},
		{"preserve", CascadePreserve},
		{"missing", CascadeMissing},
	}
	for _, tc := range cases {
		got, err := ParseCascadeStrategy(tc.input)
		if err != nil {
			t.Errorf("ParseCascadeStrategy(%q): unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseCascadeStrategy(%q) = %d, want %d", tc.input, got, tc.want)
		}
	}
}

func TestParseCascadeStrategy_Invalid(t *testing.T) {
	_, err := ParseCascadeStrategy("bogus")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestCascade_Overwrite(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	override := map[string]string{"B": "99", "C": "3"}
	res := Cascade([]map[string]string{base, override}, []string{"base", "override"}, CascadeOverwrite)

	expect := map[string]string{"A": "1", "B": "99", "C": "3"}
	for _, e := range res.Resolved {
		if expect[e.Key] != e.Value {
			t.Errorf("key %s: got %q, want %q", e.Key, e.Value, expect[e.Key])
		}
	}
	if len(res.Resolved) != 3 {
		t.Errorf("expected 3 resolved, got %d", len(res.Resolved))
	}
}

func TestCascade_Preserve(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	override := map[string]string{"B": "99", "C": "3"}
	res := Cascade([]map[string]string{base, override}, []string{"base", "override"}, CascadePreserve)

	for _, e := range res.Resolved {
		if e.Key == "B" && e.Value != "2" {
			t.Errorf("preserve: B should remain 2, got %q", e.Value)
		}
	}
	if len(res.Skipped) == 0 {
		t.Error("expected at least one skipped entry")
	}
}

func TestCascade_Missing(t *testing.T) {
	base := map[string]string{"A": "1"}
	extra := map[string]string{"A": "99", "B": "2"}
	res := Cascade([]map[string]string{base, extra}, []string{"base", "extra"}, CascadeMissing)

	resMap := map[string]string{}
	for _, e := range res.Resolved {
		resMap[e.Key] = e.Value
	}
	if resMap["A"] != "1" {
		t.Errorf("missing strategy: A should be 1, got %q", resMap["A"])
	}
	if resMap["B"] != "2" {
		t.Errorf("missing strategy: B should be 2, got %q", resMap["B"])
	}
}

func TestCascade_OriginTracked(t *testing.T) {
	base := map[string]string{"X": "base-val"}
	res := Cascade([]map[string]string{base}, []string{".env.base"}, CascadeOverwrite)
	if len(res.Resolved) != 1 || res.Resolved[0].Origin != ".env.base" {
		t.Errorf("expected origin .env.base, got %q", res.Resolved[0].Origin)
	}
}
