package diff

import (
	"os"
	"testing"
)

func writeTempSchema(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "schema*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoadSchema_Basic(t *testing.T) {
	path := writeTempSchema(t, "DB_HOST=database hostname\nDB_PORT\nAPP_KEY=application secret\n")
	s, err := LoadSchema(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(s.Keys))
	}
	if s.Keys["DB_HOST"] != "database hostname" {
		t.Errorf("expected description for DB_HOST, got %q", s.Keys["DB_HOST"])
	}
	if s.Keys["DB_PORT"] != "" {
		t.Errorf("expected empty description for DB_PORT, got %q", s.Keys["DB_PORT"])
	}
}

func TestLoadSchema_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempSchema(t, "# this is a comment\n\nFOO=bar\n")
	s, err := LoadSchema(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(s.Keys))
	}
}

func TestLoadSchema_NotFound(t *testing.T) {
	_, err := LoadSchema("/nonexistent/path/schema.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestValidateAgainstSchema_NoViolations(t *testing.T) {
	schema := &Schema{Keys: map[string]string{"FOO": "", "BAR": "desc"}}
	env := map[string]string{"FOO": "1", "BAR": "2", "EXTRA": "3"}
	v := ValidateAgainstSchema(schema, env)
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %v", v)
	}
}

func TestValidateAgainstSchema_WithViolations(t *testing.T) {
	schema := &Schema{Keys: map[string]string{"FOO": "foo key", "BAR": "", "BAZ": "baz key"}}
	env := map[string]string{"FOO": "1"}
	v := ValidateAgainstSchema(schema, env)
	if len(v) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(v))
	}
	if v[0].Key != "BAR" || v[1].Key != "BAZ" {
		t.Errorf("unexpected order: %v", v)
	}
	if v[1].Description != "baz key" {
		t.Errorf("expected description 'baz key', got %q", v[1].Description)
	}
}

func TestValidateAgainstSchema_Sorted(t *testing.T) {
	schema := &Schema{Keys: map[string]string{"ZEBRA": "", "ALPHA": "", "MANGO": ""}}
	env := map[string]string{}
	v := ValidateAgainstSchema(schema, env)
	if v[0].Key != "ALPHA" || v[1].Key != "MANGO" || v[2].Key != "ZEBRA" {
		t.Errorf("expected sorted violations, got %v", v)
	}
}
