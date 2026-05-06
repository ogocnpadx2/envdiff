package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"envdiff/internal/diff"
)

func TestParseWatchArgs_Defaults(t *testing.T) {
	wa, err := parseWatchArgs([]string{"a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if wa.LeftPath != "a.env" || wa.RightPath != "b.env" {
		t.Errorf("paths: got %q %q", wa.LeftPath, wa.RightPath)
	}
	if wa.Interval != 2*time.Second {
		t.Errorf("default interval: got %v", wa.Interval)
	}
	if wa.Format != diff.FormatText {
		t.Errorf("default format: got %v", wa.Format)
	}
}

func TestParseWatchArgs_CustomInterval(t *testing.T) {
	wa, err := parseWatchArgs([]string{"--interval=500ms", "a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if wa.Interval != 500*time.Millisecond {
		t.Errorf("interval: got %v", wa.Interval)
	}
}

func TestParseWatchArgs_CustomFormat(t *testing.T) {
	wa, err := parseWatchArgs([]string{"--format=json", "a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if wa.Format != diff.FormatJSON {
		t.Errorf("format: got %v", wa.Format)
	}
}

func TestParseWatchArgs_MissingPaths(t *testing.T) {
	_, err := parseWatchArgs([]string{"only-one.env"})
	if err == nil {
		t.Error("expected error for missing second path")
	}
}

func TestParseWatchArgs_InvalidInterval(t *testing.T) {
	_, err := parseWatchArgs([]string{"--interval=notaduration", "a.env", "b.env"})
	if err == nil {
		t.Error("expected error for invalid duration")
	}
}

func TestMatchFlag(t *testing.T) {
	var val string
	if !matchFlag("--format=json", "--format", &val) || val != "json" {
		t.Errorf("matchFlag failed: val=%q", val)
	}
	if matchFlag("--format", "--format", &val) {
		t.Error("matchFlag should not match flag without value")
	}
}

func TestRunWatch_InvalidArgs(t *testing.T) {
	var buf bytes.Buffer
	err := RunWatch([]string{}, &buf)
	if err == nil {
		t.Error("expected error with no args")
	}
	if !strings.Contains(err.Error(), "watch requires") {
		t.Errorf("unexpected error: %v", err)
	}
}
