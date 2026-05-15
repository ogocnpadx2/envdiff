package diff

import (
	"testing"
)

func TestProfileEnv_InfersTypes(t *testing.T) {
	env := map[string]string{
		"BOOL_KEY":   "true",
		"INT_KEY":    "42",
		"FLOAT_KEY":  "3.14",
		"URL_KEY":    "https://example.com",
		"PATH_KEY":   "/var/log/app",
		"STR_KEY":    "hello",
		"EMPTY_KEY":  "",
	}

	result := ProfileEnv(env)

	if len(result.Entries) != 7 {
		t.Fatalf("expected 7 entries, got %d", len(result.Entries))
	}

	byKey := make(map[string]ProfileEntry)
	for _, e := range result.Entries {
		byKey[e.Key] = e
	}

	cases := []struct {
		key      string
		wantType string
		wantNonEmpty bool
	}{
		{"BOOL_KEY", "bool", true},
		{"INT_KEY", "int", true},
		{"FLOAT_KEY", "float", true},
		{"URL_KEY", "url", true},
		{"PATH_KEY", "path", true},
		{"STR_KEY", "string", true},
		{"EMPTY_KEY", "empty", false},
	}

	for _, tc := range cases {
		e, ok := byKey[tc.key]
		if !ok {
			t.Errorf("key %s not found in profile", tc.key)
			continue
		}
		if e.Type != tc.wantType {
			t.Errorf("%s: expected type %q, got %q", tc.key, tc.wantType, e.Type)
		}
		if e.NonEmpty != tc.wantNonEmpty {
			t.Errorf("%s: expected NonEmpty=%v, got %v", tc.key, tc.wantNonEmpty, e.NonEmpty)
		}
	}
}

func TestProfileEnv_SortedEntries(t *testing.T) {
	env := map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"}
	result := ProfileEnv(env)

	if result.Entries[0].Key != "A_KEY" {
		t.Errorf("expected first key A_KEY, got %s", result.Entries[0].Key)
	}
	if result.Entries[2].Key != "Z_KEY" {
		t.Errorf("expected last key Z_KEY, got %s", result.Entries[2].Key)
	}
}

func TestProfileEnv_LengthRecorded(t *testing.T) {
	env := map[string]string{"KEY": "hello"}
	result := ProfileEnv(env)
	if result.Entries[0].Length != 5 {
		t.Errorf("expected length 5, got %d", result.Entries[0].Length)
	}
}

func TestProfileEnv_Empty(t *testing.T) {
	result := ProfileEnv(map[string]string{})
	if len(result.Entries) != 0 {
		t.Errorf("expected 0 entries for empty env")
	}
}

func TestInferType_NegativeInt(t *testing.T) {
	if inferType("-5") != "int" {
		t.Errorf("expected int for -5")
	}
}

func TestInferType_RelativePath(t *testing.T) {
	if inferType("./config") != "path" {
		t.Errorf("expected path for ./config")
	}
}
