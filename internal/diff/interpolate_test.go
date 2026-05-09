package diff

import (
	"os"
	"testing"
)

func TestInterpolateEnv_NoReferences(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	res, err := InterpolateEnv(env, DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Resolved["HOST"] != "localhost" || res.Resolved["PORT"] != "5432" {
		t.Errorf("unexpected resolved values: %v", res.Resolved)
	}
	if len(res.Unresolved) != 0 {
		t.Errorf("expected no unresolved, got %v", res.Unresolved)
	}
}

func TestInterpolateEnv_CurlyBraceStyle(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "https://${HOST}:${PORT}",
		"HOST":     "example.com",
		"PORT":     "443",
	}
	res, err := InterpolateEnv(env, DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := res.Resolved["BASE_URL"]; got != "https://example.com:443" {
		t.Errorf("BASE_URL = %q, want %q", got, "https://example.com:443")
	}
}

func TestInterpolateEnv_DollarStyle(t *testing.T) {
	env := map[string]string{
		"GREETING": "Hello $NAME",
		"NAME":     "World",
	}
	res, err := InterpolateEnv(env, DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := res.Resolved["GREETING"]; got != "Hello World" {
		t.Errorf("GREETING = %q, want %q", got, "Hello World")
	}
}

func TestInterpolateEnv_UnresolvedTracked(t *testing.T) {
	env := map[string]string{
		"DSN": "postgres://${DB_USER}:${DB_PASS}@localhost/db",
	}
	res, err := InterpolateEnv(env, DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Unresolved) != 2 {
		t.Errorf("expected 2 unresolved, got %v", res.Unresolved)
	}
}

func TestInterpolateEnv_FallbackToOS(t *testing.T) {
	os.Setenv("OS_VAR", "from-os")
	defer os.Unsetenv("OS_VAR")

	env := map[string]string{
		"VALUE": "${OS_VAR}",
	}
	opts := InterpolateOptions{FallbackToOS: true, FailOnMissing: false}
	res, err := InterpolateEnv(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := res.Resolved["VALUE"]; got != "from-os" {
		t.Errorf("VALUE = %q, want %q", got, "from-os")
	}
	if len(res.Unresolved) != 0 {
		t.Errorf("expected no unresolved, got %v", res.Unresolved)
	}
}

func TestInterpolateEnv_FailOnMissing(t *testing.T) {
	env := map[string]string{
		"KEY": "${MISSING_VAR}",
	}
	opts := InterpolateOptions{FallbackToOS: false, FailOnMissing: true}
	_, err := InterpolateEnv(env, opts)
	if err == nil {
		t.Error("expected error for missing variable, got nil")
	}
}
