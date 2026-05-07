package diff

import (
	"testing"
)

func TestParseRedactOptions_Defaults(t *testing.T) {
	opts := ParseRedactOptions("", "")
	if opts.Placeholder != DefaultPlaceholder {
		t.Errorf("expected placeholder %q, got %q", DefaultPlaceholder, opts.Placeholder)
	}
	if len(opts.Patterns) == 0 {
		t.Error("expected default patterns to be non-empty")
	}
}

func TestParseRedactOptions_Custom(t *testing.T) {
	opts := ParseRedactOptions("mysecret, token", "***")
	if opts.Placeholder != "***" {
		t.Errorf("expected placeholder ***, got %q", opts.Placeholder)
	}
	if len(opts.Patterns) != 2 {
		t.Errorf("expected 2 patterns, got %d", len(opts.Patterns))
	}
}

func TestShouldRedact_MatchesSubstring(t *testing.T) {
	opts := ParseRedactOptions("", "")
	if !shouldRedact("DB_PASSWORD", opts) {
		t.Error("expected DB_PASSWORD to be redacted")
	}
	if !shouldRedact("API_KEY", opts) {
		t.Error("expected API_KEY to be redacted")
	}
	if shouldRedact("APP_ENV", opts) {
		t.Error("expected APP_ENV not to be redacted")
	}
}

func TestRedactResult_RedactsMismatched(t *testing.T) {
	r := Result{
		Mismatched: []MismatchedKey{
			{Key: "DB_PASSWORD", LeftVal: "secret1", RightVal: "secret2"},
			{Key: "APP_ENV", LeftVal: "dev", RightVal: "prod"},
		},
		MissingInLeft:  []string{"MISSING_LEFT"},
		MissingInRight: []string{"MISSING_RIGHT"},
	}
	opts := ParseRedactOptions("", "")
	out := RedactResult(r, opts)

	if out.Mismatched[0].LeftVal != DefaultPlaceholder {
		t.Errorf("expected redacted left val, got %q", out.Mismatched[0].LeftVal)
	}
	if out.Mismatched[0].RightVal != DefaultPlaceholder {
		t.Errorf("expected redacted right val, got %q", out.Mismatched[0].RightVal)
	}
	if out.Mismatched[1].LeftVal != "dev" {
		t.Errorf("APP_ENV should not be redacted, got %q", out.Mismatched[1].LeftVal)
	}
	if len(out.MissingInLeft) != 1 || out.MissingInLeft[0] != "MISSING_LEFT" {
		t.Error("MissingInLeft should be preserved unchanged")
	}
}

func TestRedactResult_EmptyResult(t *testing.T) {
	opts := ParseRedactOptions("", "")
	out := RedactResult(Result{}, opts)
	if len(out.Mismatched) != 0 {
		t.Error("expected no mismatched entries")
	}
}
