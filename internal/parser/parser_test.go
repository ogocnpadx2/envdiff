package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseFile_Basic(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDEBUG=false\nSECRET=abc123\n")

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := EnvMap{"APP_ENV": "production", "DEBUG": "false", "SECRET": "abc123"}
	for k, v := range expected {
		if env[k] != v {
			t.Errorf("key %q: expected %q, got %q", k, v, env[k])
		}
	}
}

func TestParseFile_SkipsCommentsAndBlanks(t *testing.T) {
	content := "# this is a comment\n\nKEY=value\n"
	path := writeTempEnv(t, content)

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 {
		t.Errorf("expected 1 key, got %d", len(env))
	}
}

func TestParseFile_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `DB_URL="postgres://localhost/mydb"\nTOKEN='secret-token'\n`)

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["DB_URL"] != "postgres://localhost/mydb" {
		t.Errorf("unexpected DB_URL: %q", env["DB_URL"])
	}
}

func TestParseFile_MissingEquals(t *testing.T) {
	path := writeTempEnv(t, "BADLINE\n")

	_, err := ParseFile(path)
	if err == nil {
		t.Fatal("expected error for malformed line, got nil")
	}
}

func TestParseFile_FileNotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
