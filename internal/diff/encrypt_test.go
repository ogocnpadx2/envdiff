package diff

import (
	"strings"
	"testing"
)

func TestDefaultEncryptOptions(t *testing.T) {
	opts := DefaultEncryptOptions()
	if len(opts.HashKeys) == 0 {
		t.Fatal("expected non-empty default hash keys")
	}
	if opts.Prefix == "" {
		t.Fatal("expected non-empty prefix")
	}
}

func TestShouldHash_Matches(t *testing.T) {
	opts := DefaultEncryptOptions()
	cases := []string{"DB_PASSWORD", "API_TOKEN", "SECRET_KEY", "AUTH_CREDENTIAL"}
	for _, k := range cases {
		if !shouldHash(k, opts.HashKeys) {
			t.Errorf("expected %q to match hash keywords", k)
		}
	}
}

func TestShouldHash_NoMatch(t *testing.T) {
	opts := DefaultEncryptOptions()
	cases := []string{"DB_HOST", "PORT", "APP_ENV"}
	for _, k := range cases {
		if shouldHash(k, opts.HashKeys) {
			t.Errorf("expected %q NOT to match hash keywords", k)
		}
	}
}

func TestEncryptEnv_HashesSecrets(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "supersecret",
		"DB_HOST":     "localhost",
	}
	opts := DefaultEncryptOptions()
	out := EncryptEnv(env, opts)

	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST unchanged, got %q", out["DB_HOST"])
	}
	if !strings.HasPrefix(out["DB_PASSWORD"], opts.Prefix) {
		t.Errorf("expected DB_PASSWORD to be hashed, got %q", out["DB_PASSWORD"])
	}
	if out["DB_PASSWORD"] == "supersecret" {
		t.Error("expected DB_PASSWORD value to be replaced")
	}
}

func TestEncryptEnv_Deterministic(t *testing.T) {
	env := map[string]string{"API_TOKEN": "abc123"}
	opts := DefaultEncryptOptions()
	a := EncryptEnv(env, opts)
	b := EncryptEnv(env, opts)
	if a["API_TOKEN"] != b["API_TOKEN"] {
		t.Error("expected deterministic hashing")
	}
}

func TestEncryptEnvWithReport_TracksHashedKeys(t *testing.T) {
	env := map[string]string{
		"SECRET_KEY": "val1",
		"APP_NAME":   "envdiff",
		"DB_TOKEN":   "tok",
	}
	opts := DefaultEncryptOptions()
	_, report := EncryptEnvWithReport(env, opts)

	if len(report.HashedKeys) != 2 {
		t.Fatalf("expected 2 hashed keys, got %d: %v", len(report.HashedKeys), report.HashedKeys)
	}
	// keys should be sorted
	if report.HashedKeys[0] > report.HashedKeys[1] {
		t.Error("expected hashed keys to be sorted")
	}
}

func TestEncryptEnv_EmptyEnv(t *testing.T) {
	out := EncryptEnv(map[string]string{}, DefaultEncryptOptions())
	if len(out) != 0 {
		t.Errorf("expected empty output, got %v", out)
	}
}
