package diff

import (
	"testing"
)

func TestTransformEnv_RoundTrip_PrefixStrip(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
	}
	// Strip APP_ then re-add NEW_
	opts := TransformOptions{
		StripPrefix: "APP_",
		PrefixKeys:  "SVC_",
	}
	out := TransformEnv(env, opts)
	if _, ok := out["SVC_HOST"]; !ok {
		t.Errorf("expected SVC_HOST, got %v", out)
	}
	if _, ok := out["SVC_PORT"]; !ok {
		t.Errorf("expected SVC_PORT, got %v", out)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestTransformEnv_TrimDoesNotAlterKeys(t *testing.T) {
	env := map[string]string{"  SPACED  ": "  val  "}
	out := TransformEnv(env, TransformOptions{TrimValues: true})
	if _, ok := out["  SPACED  "]; !ok {
		t.Error("key should remain unchanged when only TrimValues is set")
	}
	if out["  SPACED  "] != "val" {
		t.Errorf("expected trimmed value 'val', got %q", out["  SPACED  "])
	}
}

func TestTransformEnv_EmptyMap(t *testing.T) {
	out := TransformEnv(map[string]string{}, TransformOptions{
		TrimValues:    true,
		UppercaseKeys: true,
		PrefixKeys:    "X_",
	})
	if len(out) != 0 {
		t.Errorf("expected empty output map, got %v", out)
	}
}

func TestTransformEnv_LowercasePreservesValues(t *testing.T) {
	env := map[string]string{"MY_KEY": "SomeValue"}
	out := TransformEnv(env, TransformOptions{LowercaseKeys: true})
	if out["my_key"] != "SomeValue" {
		t.Errorf("value should not be changed by LowercaseKeys, got %q", out["my_key"])
	}
}
