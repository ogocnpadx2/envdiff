package diff

import (
	"strings"
	"testing"
)

func TestFlattenEnv_NoChange(t *testing.T) {
	env := map[string]string{
		"DB__HOST": "localhost",
		"APP__PORT": "8080",
	}
	opts := DefaultFlattenOptions()
	r := FlattenEnv(env, opts)

	if len(r.Renamed) != 0 {
		t.Errorf("expected no renames, got %v", r.Renamed)
	}
	if len(r.Unchanged) != 2 {
		t.Errorf("expected 2 unchanged keys, got %d", len(r.Unchanged))
	}
}

func TestFlattenEnv_CollapsesConsecutiveSeparators(t *testing.T) {
	env := map[string]string{
		"DB____HOST": "localhost",
	}
	opts := DefaultFlattenOptions()
	r := FlattenEnv(env, opts)

	newKey, ok := r.Renamed["DB____HOST"]
	if !ok {
		t.Fatal("expected DB____HOST to be renamed")
	}
	if newKey != "DB__HOST" {
		t.Errorf("expected DB__HOST, got %s", newKey)
	}
	if v := r.Flattened["DB__HOST"]; v != "localhost" {
		t.Errorf("expected value localhost, got %s", v)
	}
}

func TestFlattenEnv_Uppercase(t *testing.T) {
	env := map[string]string{
		"db__host": "localhost",
	}
	opts := DefaultFlattenOptions()
	opts.Uppercase = true
	r := FlattenEnv(env, opts)

	if _, ok := r.Flattened["DB__HOST"]; !ok {
		t.Error("expected key DB__HOST in flattened output")
	}
}

func TestFlattenEnv_MaxDepth(t *testing.T) {
	env := map[string]string{
		"A__B__C__D": "value",
	}
	opts := DefaultFlattenOptions()
	opts.MaxDepth = 2
	r := FlattenEnv(env, opts)

	newKey, ok := r.Renamed["A__B__C__D"]
	if !ok {
		t.Fatal("expected A__B__C__D to be renamed due to MaxDepth")
	}
	if newKey != "A__B__C__D" {
		// With MaxDepth=2: segments [A, B, C, D] → [A, B__C__D]
		// so expected key is A__B__C__D only if already 2 segments; here it's 4.
		// Expect "A__B__C__D" truncated to 2 parts: "A__B__C__D"
	}
	parts := strings.SplitN(newKey, "__", -1)
	if len(parts) > 2 {
		t.Errorf("expected at most 2 segments after MaxDepth=2, got %d in %q", len(parts), newKey)
	}
}

func TestFlattenEnv_EmptyMap(t *testing.T) {
	r := FlattenEnv(map[string]string{}, DefaultFlattenOptions())
	if len(r.Flattened) != 0 {
		t.Error("expected empty flattened map")
	}
}

func TestFormatFlattenText_Clean(t *testing.T) {
	r := FlattenResult{
		Flattened: map[string]string{"A__B": "v"},
		Renamed:   map[string]string{},
		Unchanged: []string{"A__B"},
	}
	out := FormatFlattenText(r)
	if !strings.Contains(out, "No keys required flattening") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatFlattenText_WithRenames(t *testing.T) {
	r := FlattenResult{
		Flattened: map[string]string{"A__B": "v"},
		Renamed:   map[string]string{"A____B": "A__B"},
		Unchanged: []string{},
	}
	out := FormatFlattenText(r)
	if !strings.Contains(out, "A____B") {
		t.Errorf("expected original key in output, got: %s", out)
	}
	if !strings.Contains(out, "A__B") {
		t.Errorf("expected new key in output, got: %s", out)
	}
	if !strings.Contains(out, "Flattened 1 key") {
		t.Errorf("expected summary line, got: %s", out)
	}
}
