package diff

import (
	"strings"
	"testing"
)

func TestDefaultTruncateOptions(t *testing.T) {
	opts := DefaultTruncateOptions()
	if opts.MaxLength != 40 {
		t.Errorf("expected MaxLength 40, got %d", opts.MaxLength)
	}
	if opts.Suffix != "..." {
		t.Errorf("expected Suffix '...', got %q", opts.Suffix)
	}
}

func TestTruncateEnv_ShortValuesUnchanged(t *testing.T) {
	env := map[string]string{"KEY": "short"}
	out := TruncateEnv(env, DefaultTruncateOptions())
	if out["KEY"] != "short" {
		t.Errorf("expected 'short', got %q", out["KEY"])
	}
}

func TestTruncateEnv_LongValueTruncated(t *testing.T) {
	long := strings.Repeat("x", 80)
	env := map[string]string{"KEY": long}
	opts := DefaultTruncateOptions()
	out := TruncateEnv(env, opts)
	if len(out["KEY"]) != opts.MaxLength+len(opts.Suffix) {
		t.Errorf("unexpected length: %d", len(out["KEY"]))
	}
	if !strings.HasSuffix(out["KEY"], "...") {
		t.Errorf("expected suffix '...', got %q", out["KEY"])
	}
}

func TestTruncateEnv_CustomSuffix(t *testing.T) {
	env := map[string]string{"KEY": strings.Repeat("a", 50)}
	opts := TruncateOptions{MaxLength: 10, Suffix: "~~"}
	out := TruncateEnv(env, opts)
	if !strings.HasSuffix(out["KEY"], "~~") {
		t.Errorf("expected custom suffix, got %q", out["KEY"])
	}
}

func TestTruncateEnv_KeyFilter(t *testing.T) {
	long := strings.Repeat("z", 80)
	env := map[string]string{
		"SECRET_TOKEN": long,
		"HOST":         long,
	}
	opts := TruncateOptions{MaxLength: 10, Suffix: "...", KeyFilter: "SECRET"}
	out := TruncateEnv(env, opts)
	if len(out["SECRET_TOKEN"]) != 13 {
		t.Errorf("SECRET_TOKEN should be truncated, got len %d", len(out["SECRET_TOKEN"]))
	}
	if out["HOST"] != long {
		t.Errorf("HOST should be unchanged")
	}
}

func TestTruncateEnvWithReport_ReportsTruncated(t *testing.T) {
	env := map[string]string{
		"A": strings.Repeat("a", 60),
		"B": "short",
		"C": strings.Repeat("c", 50),
	}
	opts := DefaultTruncateOptions()
	_, truncated := TruncateEnvWithReport(env, opts)
	if len(truncated) != 2 {
		t.Errorf("expected 2 truncated keys, got %d: %v", len(truncated), truncated)
	}
}

func TestTruncateEnvWithReport_SortedKeys(t *testing.T) {
	env := map[string]string{
		"Z_KEY": strings.Repeat("z", 60),
		"A_KEY": strings.Repeat("a", 60),
	}
	opts := DefaultTruncateOptions()
	_, truncated := TruncateEnvWithReport(env, opts)
	if len(truncated) == 2 && truncated[0] != "A_KEY" {
		t.Errorf("expected sorted output, got %v", truncated)
	}
}

func TestTruncateEnv_ZeroMaxLengthUsesDefault(t *testing.T) {
	env := map[string]string{"K": strings.Repeat("x", 80)}
	opts := TruncateOptions{MaxLength: 0, Suffix: "..."}
	out := TruncateEnv(env, opts)
	def := DefaultTruncateOptions()
	expected := strings.Repeat("x", def.MaxLength) + "..."
	if out["K"] != expected {
		t.Errorf("expected default truncation, got %q", out["K"])
	}
}
