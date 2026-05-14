package diff

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempIgnoreFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".envignore")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestParseIgnoreFile_Basic(t *testing.T) {
	p := writeTempIgnoreFile(t, "SECRET_KEY\nDEBUG\n")
	opts, err := ParseIgnoreFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.Keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(opts.Keys))
	}
}

func TestParseIgnoreFile_SkipsCommentsAndBlanks(t *testing.T) {
	p := writeTempIgnoreFile(t, "# this is a comment\n\nSECRET_KEY\n")
	opts, err := ParseIgnoreFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.Keys) != 1 || opts.Keys[0] != "SECRET_KEY" {
		t.Fatalf("unexpected keys: %v", opts.Keys)
	}
}

func TestParseIgnoreFile_PrefixPattern(t *testing.T) {
	p := writeTempIgnoreFile(t, "AWS_*\nDB_*\n")
	opts, err := ParseIgnoreFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.Prefixes) != 2 {
		t.Fatalf("expected 2 prefixes, got %d", len(opts.Prefixes))
	}
}

func TestParseIgnoreFile_NotFound(t *testing.T) {
	_, err := ParseIgnoreFile("/nonexistent/.envignore")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestApplyIgnore_RemovesExactKeys(t *testing.T) {
	r := Result{
		MissingInRight: []string{"SECRET_KEY", "PORT"},
		MissingInLeft:  []string{"DEBUG"},
		Mismatched:     []Mismatch{{Key: "HOST", Left: "a", Right: "b"}},
	}
	opts := IgnoreOptions{Keys: []string{"SECRET_KEY", "DEBUG", "HOST"}}
	out := ApplyIgnore(r, opts)
	if len(out.MissingInRight) != 1 || out.MissingInRight[0] != "PORT" {
		t.Fatalf("unexpected MissingInRight: %v", out.MissingInRight)
	}
	if len(out.MissingInLeft) != 0 {
		t.Fatalf("expected empty MissingInLeft, got %v", out.MissingInLeft)
	}
	if len(out.Mismatched) != 0 {
		t.Fatalf("expected empty Mismatched, got %v", out.Mismatched)
	}
}

func TestApplyIgnore_RemovesByPrefix(t *testing.T) {
	r := Result{
		MissingInRight: []string{"AWS_KEY", "AWS_SECRET", "PORT"},
	}
	opts := IgnoreOptions{Prefixes: []string{"AWS_"}}
	out := ApplyIgnore(r, opts)
	if len(out.MissingInRight) != 1 || out.MissingInRight[0] != "PORT" {
		t.Fatalf("unexpected result: %v", out.MissingInRight)
	}
}

func TestApplyIgnore_NoOptions_ReturnsOriginal(t *testing.T) {
	r := Result{
		MissingInRight: []string{"A", "B"},
	}
	out := ApplyIgnore(r, IgnoreOptions{})
	if len(out.MissingInRight) != 2 {
		t.Fatalf("expected unchanged result, got %v", out.MissingInRight)
	}
}
