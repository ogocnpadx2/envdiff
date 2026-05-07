package diff

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func cleanResult() Result {
	return Result{}
}

func diffResult() Result {
	return Result{
		MissingInLeft:  []string{"ONLY_RIGHT"},
		MissingInRight: []string{"ONLY_LEFT"},
		Mismatched:     []MismatchedKey{{Key: "HOST", LeftValue: "localhost", RightValue: "prod.example.com"}},
	}
}

func TestRenderTemplate_Clean(t *testing.T) {
	tmpl := `clean={{ .IsClean }} diffs={{ .TotalDiffs }}`
	var buf bytes.Buffer
	if err := RenderTemplate(&buf, tmpl, cleanResult(), "a.env", "b.env"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if got != "clean=true diffs=0" {
		t.Errorf("got %q", got)
	}
}

func TestRenderTemplate_WithDiffs(t *testing.T) {
	tmpl := `left={{ .LeftFile }} right={{ .RightFile }} diffs={{ .TotalDiffs }}`
	var buf bytes.Buffer
	if err := RenderTemplate(&buf, tmpl, diffResult(), "dev.env", "prod.env"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "left=dev.env") || !strings.Contains(got, "diffs=3") {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestRenderTemplate_InvalidTemplate(t *testing.T) {
	var buf bytes.Buffer
	err := RenderTemplate(&buf, `{{ .Unclosed`, cleanResult(), "a.env", "b.env")
	if err == nil {
		t.Fatal("expected error for invalid template")
	}
}

func TestRenderTemplateFile_Basic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tmpl.txt")
	content := `files={{ .LeftFile }}+{{ .RightFile }}`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := RenderTemplateFile(&buf, path, cleanResult(), "x.env", "y.env"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got != "files=x.env+y.env" {
		t.Errorf("got %q", got)
	}
}

func TestRenderTemplateFile_NotFound(t *testing.T) {
	var buf bytes.Buffer
	err := RenderTemplateFile(&buf, "/nonexistent/tmpl.txt", cleanResult(), "a.env", "b.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
