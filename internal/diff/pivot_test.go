package diff

import (
	"strings"
	"testing"
)

func TestBuildPivot_BasicKeys(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"HOST": "localhost", "PORT": "3000"},
		"prod": {"HOST": "example.com", "PORT": "443"},
	}
	pt := BuildPivot(envs)

	if len(pt.Envs) != 2 {
		t.Fatalf("expected 2 envs, got %d", len(pt.Envs))
	}
	if pt.Envs[0] != "dev" || pt.Envs[1] != "prod" {
		t.Errorf("unexpected env order: %v", pt.Envs)
	}
	if len(pt.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(pt.Entries))
	}
	if pt.Entries[0].Key != "HOST" {
		t.Errorf("expected HOST first, got %s", pt.Entries[0].Key)
	}
}

func TestBuildPivot_MissingKey(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"HOST": "localhost", "SECRET": "abc"},
		"prod": {"HOST": "example.com"},
	}
	pt := BuildPivot(envs)

	var secretEntry *PivotEntry
	for i := range pt.Entries {
		if pt.Entries[i].Key == "SECRET" {
			secretEntry = &pt.Entries[i]
		}
	}
	if secretEntry == nil {
		t.Fatal("expected SECRET entry")
	}
	if secretEntry.Values["prod"] != "" {
		t.Errorf("expected empty string for missing prod SECRET, got %q", secretEntry.Values["prod"])
	}
}

func TestBuildPivot_Empty(t *testing.T) {
	pt := BuildPivot(map[string]map[string]string{})
	if len(pt.Entries) != 0 {
		t.Errorf("expected no entries for empty input")
	}
}

func TestFormatPivotText_ContainsHeaders(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"DB_URL": "postgres://localhost"},
		"prod": {"DB_URL": "postgres://prod"},
	}
	pt := BuildPivot(envs)
	out := FormatPivotText(pt)

	if !strings.Contains(out, "KEY") {
		t.Error("expected KEY header in output")
	}
	if !strings.Contains(out, "dev") {
		t.Error("expected dev column header")
	}
	if !strings.Contains(out, "DB_URL") {
		t.Error("expected DB_URL in output")
	}
}

func TestFormatPivotText_ShowsMissing(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"ONLY_DEV": "yes"},
		"prod": {},
	}
	pt := BuildPivot(envs)
	out := FormatPivotText(pt)

	if !strings.Contains(out, "<missing>") {
		t.Error("expected <missing> marker for absent key")
	}
}

func TestFormatPivotText_Empty(t *testing.T) {
	pt := PivotTable{}
	out := FormatPivotText(pt)
	if !strings.Contains(out, "no keys") {
		t.Errorf("expected 'no keys' message, got: %s", out)
	}
}
