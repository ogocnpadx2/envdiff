package diff

import (
	"strings"
	"testing"
)

func TestParseExportFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected ExportFormat
	}{
		{"csv", ExportCSV},
		{"CSV", ExportCSV},
		{"dotenv", ExportDotEnv},
		{"DotEnv", ExportDotEnv},
	}
	for _, tc := range cases {
		f, err := ParseExportFormat(tc.input)
		if err != nil {
			t.Errorf("ParseExportFormat(%q) unexpected error: %v", tc.input, err)
		}
		if f != tc.expected {
			t.Errorf("ParseExportFormat(%q) = %q, want %q", tc.input, f, tc.expected)
		}
	}
}

func TestParseExportFormat_Invalid(t *testing.T) {
	_, err := ParseExportFormat("xml")
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestExport_CSV(t *testing.T) {
	result := Result{
		MissingInRight: []string{"HOST"},
		MissingInLeft:  []string{"PORT"},
		Mismatched:     []MismatchedKey{{Key: "DB", LeftVal: "dev", RightVal: "prod"}},
	}
	var buf strings.Builder
	if err := Export(&buf, result, ExportCSV); err != nil {
		t.Fatalf("Export CSV error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "key,status,left_value,right_value") {
		t.Error("missing CSV header")
	}
	if !strings.Contains(out, "HOST,missing_in_right") {
		t.Error("missing HOST row")
	}
	if !strings.Contains(out, "PORT,missing_in_left") {
		t.Error("missing PORT row")
	}
	if !strings.Contains(out, "DB,mismatched,dev,prod") {
		t.Error("missing DB mismatch row")
	}
}

func TestExport_DotEnv(t *testing.T) {
	result := Result{
		MissingInRight: []string{"HOST"},
		Mismatched:     []MismatchedKey{{Key: "DB", LeftVal: "dev", RightVal: "prod"}},
	}
	var buf strings.Builder
	if err := Export(&buf, result, ExportDotEnv); err != nil {
		t.Fatalf("Export DotEnv error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "# HOST=") {
		t.Error("expected commented HOST key")
	}
	if !strings.Contains(out, "DB=dev") {
		t.Error("expected DB=dev line")
	}
}

func TestExport_Clean(t *testing.T) {
	result := Result{}
	var buf strings.Builder
	if err := Export(&buf, result, ExportDotEnv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for clean result, got: %q", buf.String())
	}
}
