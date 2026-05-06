package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func writeSnapEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeSnapEnv: %v", err)
	}
	return p
}

func TestParseSnapshotArgs_MissingPaths(t *testing.T) {
	_, err := parseSnapshotArgs([]string{})
	if err == nil {
		t.Fatal("expected error for missing paths")
	}
}

func TestParseSnapshotArgs_MissingFlag(t *testing.T) {
	_, err := parseSnapshotArgs([]string{"a.env", "b.env"})
	if err == nil {
		t.Fatal("expected error when neither --save nor --compare given")
	}
}

func TestParseSnapshotArgs_SaveFlag(t *testing.T) {
	a, err := parseSnapshotArgs([]string{"a.env", "b.env", "--save", "out.json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.outputPath != "out.json" {
		t.Errorf("outputPath: got %q", a.outputPath)
	}
}

func TestRunSnapshot_Save(t *testing.T) {
	dir := t.TempDir()
	left := writeSnapEnv(t, dir, "left.env", "FOO=1\nBAR=2\n")
	right := writeSnapEnv(t, dir, "right.env", "FOO=1\n")
	snapPath := filepath.Join(dir, "snap.json")

	var out, errBuf bytes.Buffer
	code := RunSnapshot([]string{left, right, "--save", snapPath}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, errBuf.String())
	}
	if _, err := os.Stat(snapPath); err != nil {
		t.Fatalf("snapshot file not created: %v", err)
	}
}

func TestRunSnapshot_CompareNoChange(t *testing.T) {
	dir := t.TempDir()
	left := writeSnapEnv(t, dir, "left.env", "FOO=1\n")
	right := writeSnapEnv(t, dir, "right.env", "FOO=1\n")
	snapPath := filepath.Join(dir, "snap.json")

	var out, errBuf bytes.Buffer
	// save first
	RunSnapshot([]string{left, right, "--save", snapPath}, &out, &errBuf)
	out.Reset()

	code := RunSnapshot([]string{left, right, "--compare", snapPath}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if out.String() != "no changes since snapshot\n" {
		t.Errorf("unexpected output: %q", out.String())
	}
}

func TestRunSnapshot_CompareWithNewIssue(t *testing.T) {
	dir := t.TempDir()
	left := writeSnapEnv(t, dir, "left.env", "FOO=1\n")
	right := writeSnapEnv(t, dir, "right.env", "FOO=1\n")
	snapPath := filepath.Join(dir, "snap.json")

	var out, errBuf bytes.Buffer
	RunSnapshot([]string{left, right, "--save", snapPath}, &out, &errBuf)
	out.Reset()

	// introduce a new missing key
	writeSnapEnv(t, dir, "left.env", "FOO=1\nNEW_KEY=x\n")
	code := RunSnapshot([]string{left, right, "--compare", snapPath}, &out, &errBuf)
	if code != 1 {
		t.Fatalf("expected exit 1 for new issues, got %d", code)
	}
}
