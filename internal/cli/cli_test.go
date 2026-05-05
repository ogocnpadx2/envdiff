package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestRun_Clean(t *testing.T) {
	a := writeTempEnv(t, "KEY=value\nFOO=bar\n")
	b := writeTempEnv(t, "KEY=value\nFOO=bar\n")

	var out, errOut bytes.Buffer
	err := Run([]string{a, b}, &out, &errOut)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out.String(), "No differences") {
		t.Errorf("expected clean message, got: %q", out.String())
	}
}

func TestRun_Differences(t *testing.T) {
	a := writeTempEnv(t, "KEY=value\nONLY_LEFT=yes\n")
	b := writeTempEnv(t, "KEY=other\n")

	var out, errOut bytes.Buffer
	err := Run([]string{a, b}, &out, &errOut)
	if err == nil {
		t.Fatal("expected error due to differences")
	}
	output := out.String()
	if !strings.Contains(output, "ONLY_LEFT") {
		t.Errorf("expected ONLY_LEFT in output, got: %q", output)
	}
	if !strings.Contains(output, "KEY") {
		t.Errorf("expected KEY mismatch in output, got: %q", output)
	}
}

func TestRun_QuietClean(t *testing.T) {
	a := writeTempEnv(t, "KEY=value\n")
	b := writeTempEnv(t, "KEY=value\n")

	var out, errOut bytes.Buffer
	err := Run([]string{"-quiet", a, b}, &out, &errOut)
	if err != nil {
		t.Fatalf("expected no error in quiet+clean mode, got: %v", err)
	}
	if out.Len() != 0 {
		t.Errorf("expected no output in quiet mode, got: %q", out.String())
	}
}

func TestRun_MissingArgs(t *testing.T) {
	var out, errOut bytes.Buffer
	err := Run([]string{}, &out, &errOut)
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRun_BadFile(t *testing.T) {
	var out, errOut bytes.Buffer
	err := Run([]string{"/nonexistent/.env", "/also/missing/.env"}, &out, &errOut)
	if err == nil {
		t.Fatal("expected error for missing files")
	}
}
