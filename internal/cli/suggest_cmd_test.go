package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempSuggestEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestParseSuggestArgs_MissingPaths(t *testing.T) {
	_, err := parseSuggestArgs([]string{})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestParseSuggestArgs_Defaults(t *testing.T) {
	sa, err := parseSuggestArgs([]string{"a.env", "b.env"})
	if err != nil {
		t.Fatal(err)
	}
	if sa.maxDistance != 3 {
		t.Errorf("expected default maxDistance 3, got %d", sa.maxDistance)
	}
}

func TestParseSuggestArgs_CustomDistance(t *testing.T) {
	sa, err := parseSuggestArgs([]string{"a.env", "b.env", "--max-distance=5"})
	if err != nil {
		t.Fatal(err)
	}
	if sa.maxDistance != 5 {
		t.Errorf("expected 5, got %d", sa.maxDistance)
	}
}

func TestParseSuggestArgs_InvalidDistance(t *testing.T) {
	_, err := parseSuggestArgs([]string{"a.env", "b.env", "--max-distance=0"})
	if err == nil {
		t.Fatal("expected error for distance=0")
	}
}

func TestRunSuggest_NoSuggestions(t *testing.T) {
	left := writeTempSuggestEnv(t, "FOO=1\nBAR=2\n")
	right := writeTempSuggestEnv(t, "FOO=1\nBAR=2\n")
	var buf bytes.Buffer
	if err := RunSuggest([]string{left, right}, &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No suggestions") {
		t.Errorf("expected no-suggestions message, got: %s", buf.String())
	}
}

func TestRunSuggest_TypoDetected(t *testing.T) {
	left := writeTempSuggestEnv(t, "DATABASE_URL=postgres://left\n")
	right := writeTempSuggestEnv(t, "DATABSE_URL=postgres://right\n")
	var buf bytes.Buffer
	if err := RunSuggest([]string{left, right}, &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "DATABASE_URL") {
		t.Errorf("expected DATABASE_URL in output, got: %s", out)
	}
	if !strings.Contains(out, "DATABSE_URL") {
		t.Errorf("expected DATABSE_URL suggestion in output, got: %s", out)
	}
}

func TestRunSuggest_MissingFile(t *testing.T) {
	var buf bytes.Buffer
	err := RunSuggest([]string{"/nonexistent/a.env", "/nonexistent/b.env"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
