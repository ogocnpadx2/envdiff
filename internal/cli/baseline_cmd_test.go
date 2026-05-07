package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func writeTempBaselineEnv(t *testing.T, content string) string {
	t.Helper()
	f := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return f
}

func TestParseBaselineArgs_MissingPaths(t *testing.T) {
	_, err := parseBaselineArgs([]string{"only-one"})
	if err == nil {
		t.Fatal("expected error for missing paths")
	}
}

func TestParseBaselineArgs_Defaults(t *testing.T) {
	a, err := parseBaselineArgs([]string{"a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.baseline != ".envdiff-baseline.json" {
		t.Errorf("expected default baseline path, got %q", a.baseline)
	}
	if a.save {
		t.Error("expected save=false by default")
	}
}

func TestParseBaselineArgs_SaveFlag(t *testing.T) {
	a, err := parseBaselineArgs([]string{"a.env", "b.env", "--save", "--baseline", "custom.json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !a.save {
		t.Error("expected save=true")
	}
	if a.baseline != "custom.json" {
		t.Errorf("expected custom.json, got %q", a.baseline)
	}
}

func TestRunBaseline_Save(t *testing.T) {
	left := writeTempBaselineEnv(t, "KEY=val\nSECRET=abc\n")
	right := writeTempBaselineEnv(t, "KEY=val\n")
	blPath := filepath.Join(t.TempDir(), "bl.json")
	var out, errOut bytes.Buffer
	code := RunBaseline([]string{left, right, "--save", "--baseline", blPath}, &out, &errOut)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, errOut.String())
	}
	if _, err := os.Stat(blPath); err != nil {
		t.Fatalf("baseline file not created: %v", err)
	}
}

func TestRunBaseline_NoNewIssues(t *testing.T) {
	left := writeTempBaselineEnv(t, "KEY=val\nSECRET=abc\n")
	right := writeTempBaselineEnv(t, "KEY=val\n")
	blPath := filepath.Join(t.TempDir(), "bl.json")
	var out, errOut bytes.Buffer
	// Save baseline first
	RunBaseline([]string{left, right, "--save", "--baseline", blPath}, &out, &errOut)
	out.Reset()
	// Compare against same state — no new issues
	code := RunBaseline([]string{left, right, "--baseline", blPath}, &out, &errOut)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, errOut.String())
	}
	if got := out.String(); got == "" {
		t.Error("expected output message")
	}
}

func TestRunBaseline_NewIssueDetected(t *testing.T) {
	left := writeTempBaselineEnv(t, "KEY=val\n")
	right := writeTempBaselineEnv(t, "KEY=val\n")
	blPath := filepath.Join(t.TempDir(), "bl.json")
	var out, errOut bytes.Buffer
	// Save clean baseline
	RunBaseline([]string{left, right, "--save", "--baseline", blPath}, &out, &errOut)
	out.Reset()
	// Now introduce a new difference in left
	os.WriteFile(left, []byte("KEY=val\nNEW_KEY=secret\n"), 0o644)
	code := RunBaseline([]string{left, right, "--baseline", blPath}, &out, &errOut)
	if code != 1 {
		t.Fatalf("expected exit 1 for new issues, got %d", code)
	}
}
