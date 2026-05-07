package diff

import (
	"testing"
)

func baseRenameResult() Result {
	return Result{
		MissingInRight: []string{"DB_HOST", "OLD_SECRET"},
		MissingInLeft:  []string{"DATABASE_HOST", "NEW_SECRET"},
		Mismatched:     []Mismatch{},
	}
}

func TestApplyRenames_Applied(t *testing.T) {
	r := baseRenameResult()
	renames := RenameMap{
		"DB_HOST":    "DATABASE_HOST",
		"OLD_SECRET": "NEW_SECRET",
	}

	updated, rr := ApplyRenames(r, renames)

	if len(rr.Applied) != 2 {
		t.Fatalf("expected 2 applied renames, got %d", len(rr.Applied))
	}
	if len(updated.MissingInRight) != 0 {
		t.Errorf("expected MissingInRight to be empty, got %v", updated.MissingInRight)
	}
	if len(updated.MissingInLeft) != 0 {
		t.Errorf("expected MissingInLeft to be empty, got %v", updated.MissingInLeft)
	}
}

func TestApplyRenames_Skipped(t *testing.T) {
	r := Result{
		MissingInRight: []string{"OLD_KEY"},
		MissingInLeft:  []string{"UNRELATED"},
	}
	renames := RenameMap{
		"OLD_KEY": "COMPLETELY_DIFFERENT",
	}

	_, rr := ApplyRenames(r, renames)

	if len(rr.Skipped) != 1 {
		t.Fatalf("expected 1 skipped rename, got %d", len(rr.Skipped))
	}
	if rr.Skipped[0].OldKey != "OLD_KEY" {
		t.Errorf("unexpected skipped key: %s", rr.Skipped[0].OldKey)
	}
}

func TestApplyRenames_Conflict(t *testing.T) {
	r := Result{
		MissingInRight: []string{},
		MissingInLeft:  []string{},
	}
	renames := RenameMap{
		"FOO": "BAR",
	}

	_, rr := ApplyRenames(r, renames)

	if len(rr.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(rr.Conflicts))
	}
}

func TestApplyRenames_NoRenames(t *testing.T) {
	r := baseRenameResult()
	updated, rr := ApplyRenames(r, RenameMap{})

	if len(rr.Applied) != 0 {
		t.Errorf("expected no applied renames")
	}
	if len(updated.MissingInRight) != 2 {
		t.Errorf("expected MissingInRight unchanged, got %v", updated.MissingInRight)
	}
}

func TestApplyRenames_SortedOutput(t *testing.T) {
	r := Result{
		MissingInRight: []string{"Z_OLD", "A_OLD"},
		MissingInLeft:  []string{"Z_NEW", "A_NEW"},
	}
	renames := RenameMap{
		"Z_OLD": "Z_NEW",
		"A_OLD": "A_NEW",
	}
	_, rr := ApplyRenames(r, renames)
	if len(rr.Applied) != 2 {
		t.Fatalf("expected 2 applied")
	}
	if rr.Applied[0].OldKey != "A_OLD" {
		t.Errorf("expected sorted applied renames, got first=%s", rr.Applied[0].OldKey)
	}
}
