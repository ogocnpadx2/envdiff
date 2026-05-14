package diff

import (
	"strings"
	"testing"
)

func TestBuildPatch_AddAndRemove(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"A": "1", "C": "3"}

	entries := BuildPatch(src, dst)

	opsMap := map[string]PatchOp{}
	for _, e := range entries {
		opsMap[e.Key] = e.Op
	}

	if opsMap["B"] != PatchRemove {
		t.Errorf("expected B to be removed, got %v", opsMap["B"])
	}
	if opsMap["C"] != PatchAdd {
		t.Errorf("expected C to be added, got %v", opsMap["C"])
	}
	if _, ok := opsMap["A"]; ok {
		t.Error("A should not appear in patch (unchanged)")
	}
}

func TestBuildPatch_Change(t *testing.T) {
	src := map[string]string{"HOST": "localhost"}
	dst := map[string]string{"HOST": "prod.example.com"}

	entries := BuildPatch(src, dst)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Op != PatchChange {
		t.Errorf("expected change op, got %v", e.Op)
	}
	if e.OldValue != "localhost" || e.NewValue != "prod.example.com" {
		t.Errorf("unexpected values: %v", e)
	}
}

func TestBuildPatch_Identical(t *testing.T) {
	src := map[string]string{"X": "1"}
	entries := BuildPatch(src, src)
	if len(entries) != 0 {
		t.Errorf("expected no entries for identical maps, got %d", len(entries))
	}
}

func TestApplyPatch_Add(t *testing.T) {
	target := map[string]string{"A": "1"}
	entries := []PatchEntry{{Op: PatchAdd, Key: "B", NewValue: "2"}}
	out, res := ApplyPatch(target, entries, false)
	if out["B"] != "2" {
		t.Errorf("expected B=2, got %v", out["B"])
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied, got %d", len(res.Applied))
	}
}

func TestApplyPatch_ConflictOnAdd(t *testing.T) {
	target := map[string]string{"A": "existing"}
	entries := []PatchEntry{{Op: PatchAdd, Key: "A", NewValue: "new"}}
	out, res := ApplyPatch(target, entries, false)
	if out["A"] != "existing" {
		t.Error("conflict should not overwrite existing value")
	}
	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(res.Conflicts))
	}
}

func TestApplyPatch_DryRun(t *testing.T) {
	target := map[string]string{"A": "1"}
	entries := []PatchEntry{{Op: PatchAdd, Key: "B", NewValue: "2"}}
	out, res := ApplyPatch(target, entries, true)
	if _, ok := out["B"]; ok {
		t.Error("dry run should not modify map")
	}
	if len(res.Applied) != 1 {
		t.Errorf("dry run should still report applied, got %d", len(res.Applied))
	}
}

func TestApplyPatch_ChangeConflict(t *testing.T) {
	target := map[string]string{"HOST": "staging"}
	entries := []PatchEntry{{Op: PatchChange, Key: "HOST", OldValue: "dev", NewValue: "prod"}}
	_, res := ApplyPatch(target, entries, false)
	if len(res.Conflicts) != 1 {
		t.Errorf("expected conflict when old value mismatch, got %d", len(res.Conflicts))
	}
}

func TestFormatPatch_Output(t *testing.T) {
	entries := []PatchEntry{
		{Op: PatchAdd, Key: "NEW", NewValue: "val"},
		{Op: PatchRemove, Key: "OLD", OldValue: "gone"},
		{Op: PatchChange, Key: "HOST", OldValue: "dev", NewValue: "prod"},
	}
	out := FormatPatch(entries)
	if !strings.Contains(out, "+ NEW=val") {
		t.Errorf("missing add line in: %s", out)
	}
	if !strings.Contains(out, "- OLD=gone") {
		t.Errorf("missing remove line in: %s", out)
	}
	if !strings.Contains(out, "~ HOST: dev -> prod") {
		t.Errorf("missing change line in: %s", out)
	}
}
