package diff

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/template"
)

// TemplateData holds the data passed to a custom output template.
type TemplateData struct {
	LeftFile        string
	RightFile       string
	MissingInLeft   []string
	MissingInRight  []string
	Mismatched      []MismatchedKey
	IsClean         bool
	TotalDiffs      int
}

// buildTemplateData constructs a TemplateData from a Result and file paths.
func buildTemplateData(r Result, leftFile, rightFile string) TemplateData {
	total := len(r.MissingInLeft) + len(r.MissingInRight) + len(r.Mismatched)
	return TemplateData{
		LeftFile:       leftFile,
		RightFile:      rightFile,
		MissingInLeft:  r.MissingInLeft,
		MissingInRight: r.MissingInRight,
		Mismatched:     r.Mismatched,
		IsClean:        total == 0,
		TotalDiffs:     total,
	}
}

// RenderTemplate renders a Go text/template string with diff data and writes to w.
func RenderTemplate(w io.Writer, tmplStr string, r Result, leftFile, rightFile string) error {
	tmpl, err := template.New("envdiff").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}
	data := buildTemplateData(r, leftFile, rightFile)
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	_, err = w.Write(buf.Bytes())
	return err
}

// RenderTemplateFile reads a template from a file and renders it.
func RenderTemplateFile(w io.Writer, tmplPath string, r Result, leftFile, rightFile string) error {
	raw, err := os.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("read template file %q: %w", tmplPath, err)
	}
	return RenderTemplate(w, string(raw), r, leftFile, rightFile)
}
