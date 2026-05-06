package diff

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected OutputFormat
	}{
		{"", FormatText},
		{"text", FormatText},
		{"TEXT", FormatText},
		{"json", FormatJSON},
		{"JSON", FormatJSON},
		{"markdown", FormatMarkdown},
		{"md", FormatMarkdown},
	}
	for _, tc := range cases {
		f, err := ParseFormat(tc.input)
		if err != nil {
			t.Errorf("ParseFormat(%q): unexpected error: %v", tc.input, err)
		}
		if f != tc.expected {
			t.Errorf("ParseFormat(%q) = %v, want %v", tc.input, f, tc.expected)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := ParseFormat("xml")
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestPrintFormatted_Text_Clean(t *testing.T) {
	var buf bytes.Buffer
	result := Result{}
	if err := PrintFormatted(&buf, result, "a.env", "b.env", FormatText); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected clean message, got: %s", buf.String())
	}
}

func TestPrintFormatted_JSON_Clean(t *testing.T) {
	var buf bytes.Buffer
	result := Result{}
	if err := PrintFormatted(&buf, result, "a.env", "b.env", FormatJSON); err != nil {
		t.Fatal(err)
	}
	var out jsonReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !out.Clean {
		t.Error("expected clean=true")
	}
}

func TestPrintFormatted_JSON_WithDiffs(t *testing.T) {
	var buf bytes.Buffer
	result := Result{
		MissingInRight: []string{"FOO"},
		Mismatched:     []MismatchedKey{{Key: "BAR", LeftValue: "a", RightValue: "b"}},
	}
	if err := PrintFormatted(&buf, result, "a.env", "b.env", FormatJSON); err != nil {
		t.Fatal(err)
	}
	var out jsonReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.Clean {
		t.Error("expected clean=false")
	}
	if len(out.MissingIn.Right) != 1 || out.MissingIn.Right[0] != "FOO" {
		t.Errorf("unexpected missing right: %v", out.MissingIn.Right)
	}
	if len(out.Mismatched) != 1 || out.Mismatched[0].Key != "BAR" {
		t.Errorf("unexpected mismatched: %v", out.Mismatched)
	}
}

func TestPrintFormatted_Markdown(t *testing.T) {
	var buf bytes.Buffer
	result := Result{
		MissingInLeft: []string{"SECRET"},
	}
	if err := PrintFormatted(&buf, result, "dev.env", "prod.env", FormatMarkdown); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "## envdiff") {
		t.Error("expected markdown heading")
	}
	if !strings.Contains(out, "SECRET") {
		t.Error("expected SECRET in output")
	}
}
