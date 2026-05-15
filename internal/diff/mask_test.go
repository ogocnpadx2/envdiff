package diff

import (
	"regexp"
	"testing"
)

func TestDefaultMaskOptions(t *testing.T) {
	opts := DefaultMaskOptions()
	if opts.MaskChar != "*" {
		t.Fatalf("expected MaskChar '*', got %q", opts.MaskChar)
	}
	if opts.RevealChars != 0 {
		t.Fatalf("expected RevealChars 0, got %d", opts.RevealChars)
	}
	if len(opts.Patterns) == 0 {
		t.Fatal("expected default patterns to be non-empty")
	}
}

func TestShouldMask_MatchesPattern(t *testing.T) {
	opts := DefaultMaskOptions()
	if !shouldMask("DB_PASSWORD", opts) {
		t.Error("expected DB_PASSWORD to be masked")
	}
	if !shouldMask("api_secret", opts) {
		t.Error("expected api_secret to be masked")
	}
	if shouldMask("APP_NAME", opts) {
		t.Error("expected APP_NAME not to be masked")
	}
}

func TestShouldMask_Regex(t *testing.T) {
	opts := DefaultMaskOptions()
	opts.Regex = regexp.MustCompile(`(?i)^INTERNAL_`)
	if !shouldMask("INTERNAL_TOKEN", opts) {
		t.Error("expected INTERNAL_TOKEN to match regex")
	}
	if shouldMask("PUBLIC_URL", opts) {
		t.Error("expected PUBLIC_URL not to match regex")
	}
}

func TestMaskValue_FullMask(t *testing.T) {
	opts := DefaultMaskOptions()
	got := maskValue("supersecret", opts)
	if got != "***********" {
		t.Fatalf("expected full mask, got %q", got)
	}
}

func TestMaskValue_RevealChars(t *testing.T) {
	opts := DefaultMaskOptions()
	opts.RevealChars = 3
	got := maskValue("supersecret", opts)
	if got != "sup********" {
		t.Fatalf("expected 'sup********', got %q", got)
	}
}

func TestMaskValue_EmptyValue(t *testing.T) {
	opts := DefaultMaskOptions()
	got := maskValue("", opts)
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestMaskEnv_MasksSecrets(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "hunter2",
		"APP_NAME":    "myapp",
		"API_TOKEN":   "tok_abc123",
	}
	opts := DefaultMaskOptions()
	out := MaskEnv(env, opts)
	if out["DB_PASSWORD"] != "*******" {
		t.Errorf("expected DB_PASSWORD masked, got %q", out["DB_PASSWORD"])
	}
	if out["API_TOKEN"] != "**********" {
		t.Errorf("expected API_TOKEN masked, got %q", out["API_TOKEN"])
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME unchanged, got %q", out["APP_NAME"])
	}
}

func TestMaskEnvWithReport_ReturnsSortedKeys(t *testing.T) {
	env := map[string]string{
		"Z_SECRET": "zzz",
		"A_TOKEN":  "aaa",
		"PLAIN":    "visible",
	}
	opts := DefaultMaskOptions()
	_, masked := MaskEnvWithReport(env, opts)
	if len(masked) != 2 {
		t.Fatalf("expected 2 masked keys, got %d", len(masked))
	}
	if masked[0] != "A_TOKEN" || masked[1] != "Z_SECRET" {
		t.Errorf("expected sorted masked keys, got %v", masked)
	}
}
