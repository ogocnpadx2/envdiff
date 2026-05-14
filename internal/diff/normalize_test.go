package diff

import (
	"testing"
)

func TestNormalizeEnv_TrimSpace(t *testing.T) {
	env := map[string]string{
		"KEY": "  value  ",
		"OTHER": "clean",
	}
	opts := DefaultNormalizeOptions()
	out := NormalizeEnv(env, opts)
	if out["KEY"] != "value" {
		t.Errorf("expected trimmed value, got %q", out["KEY"])
	}
	if out["OTHER"] != "clean" {
		t.Errorf("expected unchanged value, got %q", out["OTHER"])
	}
}

func TestNormalizeEnv_LowercaseVal(t *testing.T) {
	env := map[string]string{"KEY": "Hello"}
	opts := NormalizeOptions{LowercaseVal: true}
	out := NormalizeEnv(env, opts)
	if out["KEY"] != "hello" {
		t.Errorf("expected lowercase value, got %q", out["KEY"])
	}
}

func TestNormalizeEnv_LowercaseKey(t *testing.T) {
	env := map[string]string{"MY_KEY": "val"}
	opts := NormalizeOptions{LowercaseKey: true}
	out := NormalizeEnv(env, opts)
	if _, ok := out["my_key"]; !ok {
		t.Error("expected lowercased key 'my_key' to exist")
	}
	if _, ok := out["MY_KEY"]; ok {
		t.Error("expected original key 'MY_KEY' to be absent")
	}
}

func TestNormalizeEnv_CollapseEmpty(t *testing.T) {
	env := map[string]string{
		"EMPTY": "",
		"FULL":  "value",
	}
	opts := NormalizeOptions{CollapseEmpty: true}
	out := NormalizeEnv(env, opts)
	if _, ok := out["EMPTY"]; ok {
		t.Error("expected empty key to be dropped")
	}
	if out["FULL"] != "value" {
		t.Errorf("expected FULL to be preserved, got %q", out["FULL"])
	}
}

func TestNormalizeEnvWithReport_TrimmedReported(t *testing.T) {
	env := map[string]string{
		"PADDED":  "  hello  ",
		"CLEAN":   "world",
	}
	opts := NormalizeOptions{TrimSpace: true}
	_, report := NormalizeEnvWithReport(env, opts)
	if len(report.Trimmed) != 1 || report.Trimmed[0] != "PADDED" {
		t.Errorf("expected PADDED in Trimmed, got %v", report.Trimmed)
	}
}

func TestNormalizeEnvWithReport_DroppedReported(t *testing.T) {
	env := map[string]string{
		"EMPTY": "",
		"FULL":  "v",
	}
	opts := NormalizeOptions{CollapseEmpty: true}
	_, report := NormalizeEnvWithReport(env, opts)
	if len(report.Dropped) != 1 || report.Dropped[0] != "EMPTY" {
		t.Errorf("expected EMPTY in Dropped, got %v", report.Dropped)
	}
}

func TestNormalizeEnvWithReport_RenamedReported(t *testing.T) {
	env := map[string]string{"MyKey": "val"}
	opts := NormalizeOptions{LowercaseKey: true}
	out, report := NormalizeEnvWithReport(env, opts)
	if len(report.Renamed) != 1 || report.Renamed[0] != "MyKey" {
		t.Errorf("expected MyKey in Renamed, got %v", report.Renamed)
	}
	if out["mykey"] != "val" {
		t.Errorf("expected mykey=val in output, got %v", out)
	}
}

func TestNormalizeEnvWithReport_NoChanges(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	opts := DefaultNormalizeOptions()
	_, report := NormalizeEnvWithReport(env, opts)
	if len(report.Trimmed) != 0 || len(report.Renamed) != 0 || len(report.Dropped) != 0 {
		t.Errorf("expected empty report, got %+v", report)
	}
}
