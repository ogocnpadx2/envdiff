package diff

import (
	"testing"
)

func TestSeverityString(t *testing.T) {
	cases := []struct {
		s    Severity
		want string
	}{
		{SeverityInfo, "info"},
		{SeverityWarning, "warning"},
		{SeverityCritical, "critical"},
		{Severity(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.s.String(); got != tc.want {
			t.Errorf("Severity(%d).String() = %q, want %q", tc.s, got, tc.want)
		}
	}
}

func TestClassifyResult_Clean(t *testing.T) {
	r := Result{}
	entries := ClassifyResult(r)
	if len(entries) != 0 {
		t.Errorf("expected no entries for clean result, got %d", len(entries))
	}
}

func TestClassifyResult_MissingInRight(t *testing.T) {
	r := Result{MissingInRight: []string{"DB_HOST", "DB_PORT"}}
	entries := ClassifyResult(r)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Severity != SeverityWarning {
			t.Errorf("key %q: expected Warning, got %s", e.Key, e.Severity)
		}
	}
}

func TestClassifyResult_MissingInLeft(t *testing.T) {
	r := Result{MissingInLeft: []string{"DB_HOST", "DB_PORT"}}
	entries := ClassifyResult(r)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Severity != SeverityWarning {
			t.Errorf("key %q: expected Warning, got %s", e.Key, e.Severity)
		}
	}
}

func TestClassifyResult_Mismatched(t *testing.T) {
	r := Result{
		Mismatched: []MismatchedKey{
			{Key: "API_URL", LeftVal: "http://dev", RightVal: "http://prod"},
		},
	}
	entries := ClassifyResult(r)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Severity != SeverityCritical {
		t.Errorf("expected Critical, got %s", entries[0].Severity)
	}
	if entries[0].Key != "API_URL" {
		t.Errorf("expected key API_URL, got %s", entries[0].Key)
	}
}

func TestMaxSeverity_Clean(t *testing.T) {
	r := Result{}
	if got := MaxSeverity(r); got != SeverityInfo {
		t.Errorf("expected Info for clean result, got %s", got)
	}
}

func TestMaxSeverity_MixedFindings(t *testing.T) {
	r := Result{
		MissingInLeft: []string{"SECRET"},
		Mismatched: []MismatchedKey{
			{Key: "API_URL", LeftVal: "a", RightVal: "b"},
		},
	}
	if got := MaxSeverity(r); got != SeverityCritical {
		t.Errorf("expected Critical, got %s", got)
	}
}

func TestMaxSeverity_OnlyWarning(t *testing.T) {
	r := Result{MissingInRight: []string{"OPTIONAL_KEY"}}
	if got := MaxSeverity(r); got != SeverityWarning {
		t.Errorf("expected Warning, got %s", got)
	}
}
