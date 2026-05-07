package diff

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeTempBaseline(t *testing.T, b Baseline) string {
	t.Helper()
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	f := filepath.Join(t.TempDir(), "baseline.json")
	if err := os.WriteFile(f, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return f
}

func TestSaveAndLoadBaseline(t *testing.T) {
	result := Result{
		MissingInRight: []string{"SECRET"},
		Mismatched:     []MismatchedKey{{Key: "PORT", LeftVal: "8080", RightVal: "9090"}},
	}
	path := filepath.Join(t.TempDir(), "bl.json")
	if err := SaveBaseline(path, "a.env", "b.env", result); err != nil {
		t.Fatalf("SaveBaseline: %v", err)
	}
	b, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}
	if b.LeftFile != "a.env" || b.RightFile != "b.env" {
		t.Errorf("unexpected files: %s %s", b.LeftFile, b.RightFile)
	}
	if side, ok := b.MissingKeys["SECRET"]; !ok || side != "right" {
		t.Errorf("expected SECRET missing in right, got %q", side)
	}
	if vals, ok := b.Mismatched["PORT"]; !ok || vals[0] != "8080" || vals[1] != "9090" {
		t.Errorf("unexpected PORT mismatch: %v", vals)
	}
}

func TestLoadBaseline_NotFound(t *testing.T) {
	_, err := LoadBaseline("/nonexistent/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadBaseline_InvalidJSON(t *testing.T) {
	f := filepath.Join(t.TempDir(), "bad.json")
	os.WriteFile(f, []byte("not json{"), 0o644)
	_, err := LoadBaseline(f)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestDiffAgainstBaseline_NewIssues(t *testing.T) {
	baseline := Baseline{
		MissingKeys: map[string]string{"OLD_KEY": "right"},
		Mismatched:  map[string][2]string{},
	}
	current := Result{
		MissingInRight: []string{"OLD_KEY", "NEW_KEY"},
	}
	newIssues, resolved := DiffAgainstBaseline(baseline, current)
	if len(newIssues.MissingInRight) != 1 || newIssues.MissingInRight[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY as new issue, got %v", newIssues.MissingInRight)
	}
	if len(resolved.MissingInRight) != 0 {
		t.Errorf("expected no resolved, got %v", resolved.MissingInRight)
	}
}

func TestDiffAgainstBaseline_ResolvedIssues(t *testing.T) {
	baseline := Baseline{
		MissingKeys: map[string]string{"GONE_KEY": "right"},
		Mismatched:  map[string][2]string{},
	}
	current := Result{} // GONE_KEY is now present on both sides
	_, resolved := DiffAgainstBaseline(baseline, current)
	if len(resolved.MissingInRight) != 1 || resolved.MissingInRight[0] != "GONE_KEY" {
		t.Errorf("expected GONE_KEY resolved, got %v", resolved.MissingInRight)
	}
}

func TestDiffAgainstBaseline_NewMismatch(t *testing.T) {
	baseline := Baseline{
		MissingKeys: map[string]string{},
		Mismatched:  map[string][2]string{},
	}
	current := Result{
		Mismatched: []MismatchedKey{{Key: "DB_HOST", LeftVal: "localhost", RightVal: "db.prod"}},
	}
	newIssues, _ := DiffAgainstBaseline(baseline, current)
	if len(newIssues.Mismatched) != 1 || newIssues.Mismatched[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST as new mismatch, got %v", newIssues.Mismatched)
	}
}
