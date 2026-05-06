package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempLintEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempLintEnv: %v", err)
	}
	return p
}

func TestParseLintArgs_MissingPath(t *testing.T) {
	_, err := parseLintArgs([]string{})
	if err == nil {
		t.Error("expected error for missing path")
	}
}

func TestParseLintArgs_Defaults(t *testing.T) {
	la, err := parseLintArgs([]string{"/some/file"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if la.path != "/some/file" {
		t.Errorf("unexpected path: %s", la.path)
	}
	if len(la.rules) == 0 {
		t.Error("expected default rules")
	}
}

func TestParseLintArgs_CustomRules(t *testing.T) {
	la, err := parseLintArgs([]string{"/f", "--rules=uppercase,non-empty"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(la.rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(la.rules))
	}
}

func TestParseLintArgs_InvalidRule(t *testing.T) {
	_, err := parseLintArgs([]string{"/f", "--rules=bad-rule"})
	if err == nil {
		t.Error("expected error for invalid rule")
	}
}

func TestRunLint_Clean(t *testing.T) {
	p := writeTempLintEnv(t, "APP_NAME=myapp\nPORT=8080\n")
	var out, errOut bytes.Buffer
	code := RunLint([]string{p, "--rules=uppercase,non-empty"}, &out, &errOut)
	if code != 0 {
		t.Errorf("expected exit 0, got %d; stderr: %s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "OK") {
		t.Errorf("expected OK in output, got: %s", out.String())
	}
}

func TestRunLint_WithViolations(t *testing.T) {
	p := writeTempLintEnv(t, "lower_key=\nGOOD=value\n")
	var out, errOut bytes.Buffer
	code := RunLint([]string{p, "--rules=uppercase,non-empty"}, &out, &errOut)
	if code != 1 {
		t.Errorf("expected exit 1, got %d", code)
	}
	if !strings.Contains(out.String(), "violation") {
		t.Errorf("expected violations in output, got: %s", out.String())
	}
}

func TestRunLint_MissingFile(t *testing.T) {
	var out, errOut bytes.Buffer
	code := RunLint([]string{"/nonexistent/.env"}, &out, &errOut)
	if code != 1 {
		t.Errorf("expected exit 1 for missing file, got %d", code)
	}
}
