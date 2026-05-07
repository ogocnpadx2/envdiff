package diff

import (
	"testing"
)

func TestParseTransformOptions_Valid(t *testing.T) {
	opts, err := ParseTransformOptions([]string{"trim", "uppercase-keys", "prefix=APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.TrimValues {
		t.Error("expected TrimValues=true")
	}
	if !opts.UppercaseKeys {
		t.Error("expected UppercaseKeys=true")
	}
	if opts.PrefixKeys != "APP_" {
		t.Errorf("expected PrefixKeys=APP_, got %q", opts.PrefixKeys)
	}
}

func TestParseTransformOptions_Invalid(t *testing.T) {
	_, err := ParseTransformOptions([]string{"unknown-flag"})
	if err == nil {
		t.Error("expected error for unknown flag")
	}
}

func TestParseTransformOptions_ConflictingCase(t *testing.T) {
	_, err := ParseTransformOptions([]string{"lowercase-keys", "uppercase-keys"})
	if err == nil {
		t.Error("expected error for conflicting case flags")
	}
}

func TestTransformEnv_Trim(t *testing.T) {
	env := map[string]string{"KEY": "  value  "}
	out := TransformEnv(env, TransformOptions{TrimValues: true})
	if out["KEY"] != "value" {
		t.Errorf("expected 'value', got %q", out["KEY"])
	}
}

func TestTransformEnv_LowercaseKeys(t *testing.T) {
	env := map[string]string{"MY_KEY": "v"}
	out := TransformEnv(env, TransformOptions{LowercaseKeys: true})
	if _, ok := out["my_key"]; !ok {
		t.Error("expected lowercase key 'my_key'")
	}
}

func TestTransformEnv_PrefixAndStrip(t *testing.T) {
	env := map[string]string{"OLD_KEY": "v"}
	out := TransformEnv(env, TransformOptions{StripPrefix: "OLD_", PrefixKeys: "NEW_"})
	if _, ok := out["NEW_KEY"]; !ok {
		t.Errorf("expected key 'NEW_KEY', got %v", out)
	}
}

func TestTransformEnv_UppercaseKeys(t *testing.T) {
	env := map[string]string{"lower_key": "v"}
	out := TransformEnv(env, TransformOptions{UppercaseKeys: true})
	if _, ok := out["LOWER_KEY"]; !ok {
		t.Error("expected uppercase key 'LOWER_KEY'")
	}
}

func TestTransformEnv_NoOptions(t *testing.T) {
	env := map[string]string{"KEY": "  val  "}
	out := TransformEnv(env, TransformOptions{})
	if out["KEY"] != "  val  " {
		t.Errorf("expected unchanged value, got %q", out["KEY"])
	}
}
