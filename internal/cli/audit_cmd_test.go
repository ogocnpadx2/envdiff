package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempAuditEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestParseAuditArgs_MissingPaths(t *testing.T) {
	_, err := parseAuditArgs([]string{})
	if err == nil {
		t.Fatal("expected error for missing paths")
	}
}

func TestParseAuditArgs_Defaults(t *testing.T) {
	a, err := parseAuditArgs([]string{"a.env", "b.env"})
	if err != nil {
		t.Fatal(err)
	}
	if a.leftPath != "a.env" || a.rightPath != "b.env" {
		t.Errorf("unexpected paths: %s %s", a.leftPath, a.rightPath)
	}
	if a.jsonOut {
		t.Error("jsonOut should default to false")
	}
}

func TestParseAuditArgs_JSONFlag(t *testing.T) {
	a, err := parseAuditArgs([]string{"a.env", "b.env", "--json"})
	if err != nil {
		t.Fatal(err)
	}
	if !a.jsonOut {
		t.Error("expected jsonOut to be true")
	}
}

func TestRunAudit_TextOutput(t *testing.T) {
	left := writeTempAuditEnv(t, "KEY=same\nOLD=gone\nCHANGED=before\n")
	right := writeTempAuditEnv(t, "KEY=same\nNEW=arrived\nCHANGED=after\n")

	var buf bytes.Buffer
	err := RunAudit([]string{left, right}, &buf)
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Summary:") {
		t.Error("expected Summary line in output")
	}
	if !strings.Contains(out, "+ NEW") {
		t.Error("expected added key NEW in output")
	}
	if !strings.Contains(out, "- OLD") {
		t.Error("expected removed key OLD in output")
	}
	if !strings.Contains(out, "~ CHANGED") {
		t.Error("expected changed key CHANGED in output")
	}
}

func TestRunAudit_JSONOutput(t *testing.T) {
	left := writeTempAuditEnv(t, "FOO=bar\n")
	right := writeTempAuditEnv(t, "FOO=baz\n")

	var buf bytes.Buffer
	err := RunAudit([]string{left, right, "--json"}, &buf)
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, `"type"`) {
		t.Error("expected JSON output with 'type' field")
	}
	if !strings.Contains(out, `"changed"`) {
		t.Error("expected changed event in JSON output")
	}
}

func TestRunAudit_MissingFile(t *testing.T) {
	var buf bytes.Buffer
	err := RunAudit([]string{"/no/such/file.env", "/also/missing.env"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing files")
	}
}
