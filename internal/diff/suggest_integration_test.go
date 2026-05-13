package diff

import (
	"testing"
)

// Integration tests exercise Suggest end-to-end through Compare.

func TestSuggest_Integration_RealTypo(t *testing.T) {
	left := map[string]string{
		"SECRET_KEY":   "abc",
		"DATABASE_URL": "postgres://",
	}
	right := map[string]string{
		"SECRET_KEY":  "abc",
		"DATABSE_URL": "postgres://",
	}
	result := Compare(left, right)
	sr := Suggest(result, 3)
	if len(sr.Suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(sr.Suggestions))
	}
	if sr.Suggestions[0].MissingKey != "DATABASE_URL" {
		t.Errorf("unexpected missing key: %s", sr.Suggestions[0].MissingKey)
	}
	if sr.Suggestions[0].SuggestedKey != "DATABSE_URL" {
		t.Errorf("unexpected suggestion: %s", sr.Suggestions[0].SuggestedKey)
	}
}

func TestSuggest_Integration_NoisyKeys(t *testing.T) {
	left := map[string]string{
		"ALPHA": "1",
		"BETA":  "2",
	}
	right := map[string]string{
		"GAMMA": "3",
		"DELTA": "4",
	}
	result := Compare(left, right)
	sr := Suggest(result, 2)
	// ALPHA/BETA vs GAMMA/DELTA are too different — no suggestions expected
	if len(sr.Suggestions) != 0 {
		t.Errorf("expected 0 suggestions for unrelated keys, got %d", len(sr.Suggestions))
	}
}

func TestSuggest_Integration_CaseInsensitiveMatch(t *testing.T) {
	left := map[string]string{"API_KEY": "x"}
	right := map[string]string{"api_key": "x"}
	result := Compare(left, right)
	sr := Suggest(result, 3)
	// "API_KEY" vs "api_key" differ only in case — distance should be 0 after lowercasing
	if len(sr.Suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(sr.Suggestions))
	}
	if sr.Suggestions[0].Distance != 0 {
		t.Errorf("expected distance 0 for case-only difference, got %d", sr.Suggestions[0].Distance)
	}
}

func TestSuggest_Integration_EmptyEnvs(t *testing.T) {
	left := map[string]string{}
	right := map[string]string{}
	result := Compare(left, right)
	sr := Suggest(result, 3)
	if len(sr.Suggestions) != 0 {
		t.Errorf("expected 0 suggestions for empty envs, got %d", len(sr.Suggestions))
	}
}
