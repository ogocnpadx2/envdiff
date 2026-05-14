package diff

import (
	"testing"
)

func TestParsePromoteStrategy_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  PromoteStrategy
	}{
		{"missing", PromoteOnlyMissing},
		{"", PromoteOnlyMissing},
		{"overwrite", PromoteOverwrite},
	}
	for _, tc := range cases {
		got, err := ParsePromoteStrategy(tc.input)
		if err != nil {
			t.Errorf("ParsePromoteStrategy(%q): unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParsePromoteStrategy(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestParsePromoteStrategy_Invalid(t *testing.T) {
	_, err := ParsePromoteStrategy("bogus")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestPromote_OnlyMissing(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{"A": "old", "D": "4"}

	res := Promote(src, dst, PromoteOptions{Strategy: PromoteOnlyMissing})

	if len(res.Added) != 2 { // B and C
		t.Errorf("expected 2 added, got %v", res.Added)
	}
	if len(res.Overwritten) != 0 {
		t.Errorf("expected 0 overwritten, got %v", res.Overwritten)
	}
	if len(res.Skipped) != 1 { // A already present
		t.Errorf("expected 1 skipped, got %v", res.Skipped)
	}
	if res.Merged["A"] != "old" {
		t.Errorf("A should remain old, got %q", res.Merged["A"])
	}
}

func TestPromote_Overwrite(t *testing.T) {
	src := map[string]string{"A": "new", "B": "2"}
	dst := map[string]string{"A": "old"}

	res := Promote(src, dst, PromoteOptions{Strategy: PromoteOverwrite})

	if len(res.Overwritten) != 1 {
		t.Errorf("expected 1 overwritten, got %v", res.Overwritten)
	}
	if res.Merged["A"] != "new" {
		t.Errorf("A should be overwritten to 'new', got %q", res.Merged["A"])
	}
	if res.Merged["B"] != "2" {
		t.Errorf("B should be added as '2', got %q", res.Merged["B"])
	}
}

func TestPromote_DryRun(t *testing.T) {
	src := map[string]string{"X": "1"}
	dst := map[string]string{"Y": "2"}

	res := Promote(src, dst, PromoteOptions{Strategy: PromoteOnlyMissing, DryRun: true})

	if len(res.Added) != 1 {
		t.Errorf("expected 1 added in dry-run report, got %v", res.Added)
	}
	if _, ok := res.Merged["X"]; ok {
		t.Error("dry-run should not mutate merged map")
	}
}

func TestPromote_DoesNotMutateInputs(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{"B": "2"}

	Promote(src, dst, PromoteOptions{})

	if _, ok := dst["A"]; ok {
		t.Error("Promote must not mutate the target map")
	}
}
