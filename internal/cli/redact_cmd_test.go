package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempRedactEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestParseRedactArgs_MissingPaths(t *testing.T) {
	_, err := parseRedactArgs([]string{})
	if err == nil {
		t.Error("expected error for missing paths")
	}
}

func TestParseRedactArgs_Defaults(t *testing.T) {
	ra, err := parseRedactArgs([]string{"a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ra.leftPath != "a.env" || ra.rightPath != "b.env" {
		t.Errorf("unexpected paths: %q %q", ra.leftPath, ra.rightPath)
	}
	if ra.format != "text" {
		t.Errorf("expected default format text, got %q", ra.format)
	}
}

func TestRunRedact_RedactsSecrets(t *testing.T) {
	left := writeTempRedactEnv(t, "APP_ENV=dev\nDB_PASSWORD=hunter2\n")
	right := writeTempRedactEnv(t, "APP_ENV=prod\nDB_PASSWORD=s3cr3t\n")

	var buf bytes.Buffer
	err := RunRedact([]string{left, right}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if strings.Contains(output, "hunter2") || strings.Contains(output, "s3cr3t") {
		t.Error("expected passwords to be redacted in output")
	}
	if !strings.Contains(output, "[REDACTED]") {
		t.Error("expected [REDACTED] placeholder in output")
	}
}

func TestRunRedact_PreservesNonSecretValues(t *testing.T) {
	left := writeTempRedactEnv(t, "APP_ENV=dev\nDB_PASSWORD=hunter2\n")
	right := writeTempRedactEnv(t, "APP_ENV=prod\nDB_PASSWORD=s3cr3t\n")

	var buf bytes.Buffer
	if err := RunRedact([]string{left, right}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "APP_ENV") {
		t.Error("expected APP_ENV to appear in output")
	}
}

func TestRunRedact_MissingFile(t *testing.T) {
	var buf bytes.Buffer
	err := RunRedact([]string{"/nonexistent/a.env", "/nonexistent/b.env"}, &buf)
	if err == nil {
		t.Error("expected error for missing files")
	}
}
