package diff

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveAndLoadSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	result := Result{
		MissingInRight: []string{"FOO"},
		MissingInLeft:  []string{"BAR"},
		Mismatched:     []Mismatch{{Key: "BAZ", LeftVal: "a", RightVal: "b"}},
	}

	if err := SaveSnapshot(path, "left.env", "right.env", result); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	snap, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}

	if snap.LeftFile != "left.env" {
		t.Errorf("LeftFile: got %q, want %q", snap.LeftFile, "left.env")
	}
	if snap.RightFile != "right.env" {
		t.Errorf("RightFile: got %q, want %q", snap.RightFile, "right.env")
	}
	if len(snap.Result.MissingInRight) != 1 || snap.Result.MissingInRight[0] != "FOO" {
		t.Errorf("MissingInRight mismatch: %v", snap.Result.MissingInRight)
	}
	if len(snap.Result.Mismatched) != 1 || snap.Result.Mismatched[0].Key != "BAZ" {
		t.Errorf("Mismatched mismatch: %v", snap.Result.Mismatched)
	}
}

func TestLoadSnapshot_NotFound(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadSnapshot_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0644)
	_, err := LoadSnapshot(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestDiffSnapshots_NewAndResolved(t *testing.T) {
	now := time.Now().UTC()

	before := &Snapshot{
		Timestamp: now.Add(-time.Minute),
		Result: Result{
			MissingInRight: []string{"OLD_KEY"},
		},
	}
	after := &Snapshot{
		Timestamp: now,
		Result: Result{
			MissingInRight: []string{"NEW_KEY"},
		},
	}

	delta := DiffSnapshots(before, after)

	if len(delta.NewIssues) != 1 || delta.NewIssues[0] != "NEW_KEY" {
		t.Errorf("NewIssues: got %v", delta.NewIssues)
	}
	if len(delta.ResolvedIssues) != 1 || delta.ResolvedIssues[0] != "OLD_KEY" {
		t.Errorf("ResolvedIssues: got %v", delta.ResolvedIssues)
	}
}

func TestDiffSnapshots_NoChange(t *testing.T) {
	now := time.Now().UTC()
	result := Result{MissingInRight: []string{"KEY"}}
	before := &Snapshot{Timestamp: now.Add(-time.Minute), Result: result}
	after := &Snapshot{Timestamp: now, Result: result}

	delta := DiffSnapshots(before, after)
	if len(delta.NewIssues) != 0 || len(delta.ResolvedIssues) != 0 {
		t.Errorf("expected no changes, got new=%v resolved=%v", delta.NewIssues, delta.ResolvedIssues)
	}
}
