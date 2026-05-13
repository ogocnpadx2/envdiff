package diff

import (
	"testing"
)

func TestSuggest_NoSuggestions_WhenClean(t *testing.T) {
	result := CompareResult{}
	sr := Suggest(result, 3)
	if len(sr.Suggestions) != 0 {
		t.Fatalf("expected no suggestions, got %d", len(sr.Suggestions))
	}
}

func TestSuggest_TypoDetected(t *testing.T) {
	result := CompareResult{
		MissingInRight: []string{"DATABASE_URL"},
		MissingInLeft:  []string{"DATABSE_URL"},
	}
	sr := Suggest(result, 3)
	if len(sr.Suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(sr.Suggestions))
	}
	s := sr.Suggestions[0]
	if s.MissingKey != "DATABASE_URL" {
		t.Errorf("unexpected missing key: %s", s.MissingKey)
	}
	if s.SuggestedKey != "DATABSE_URL" {
		t.Errorf("unexpected suggestion: %s", s.SuggestedKey)
	}
	if s.Distance > 3 {
		t.Errorf("distance too large: %d", s.Distance)
	}
}

func TestSuggest_NoMatchBeyondThreshold(t *testing.T) {
	result := CompareResult{
		MissingInRight: []string{"FOO"},
		MissingInLeft:  []string{"COMPLETELY_DIFFERENT_KEY"},
	}
	sr := Suggest(result, 3)
	if len(sr.Suggestions) != 0 {
		t.Fatalf("expected no suggestions, got %d", len(sr.Suggestions))
	}
}

func TestSuggest_DefaultMaxDistance(t *testing.T) {
	result := CompareResult{
		MissingInRight: []string{"API_KEY"},
		MissingInLeft:  []string{"API_KEYS"},
	}
	// passing 0 should fall back to default of 3
	sr := Suggest(result, 0)
	if len(sr.Suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(sr.Suggestions))
	}
}

func TestSuggest_SortedByDistance(t *testing.T) {
	result := CompareResult{
		MissingInRight: []string{"FOO_BAR", "FOO_BAZ"},
		MissingInLeft:  []string{"FOO_BAR_X", "FOO_BA"},
	}
	sr := Suggest(result, 5)
	for i := 1; i < len(sr.Suggestions); i++ {
		if sr.Suggestions[i].Distance < sr.Suggestions[i-1].Distance {
			t.Error("suggestions not sorted by distance")
		}
	}
}

func TestLevenshtein_SameString(t *testing.T) {
	if d := levenshtein("hello", "hello"); d != 0 {
		t.Errorf("expected 0, got %d", d)
	}
}

func TestLevenshtein_EmptyStrings(t *testing.T) {
	if d := levenshtein("", "abc"); d != 3 {
		t.Errorf("expected 3, got %d", d)
	}
	if d := levenshtein("abc", ""); d != 3 {
		t.Errorf("expected 3, got %d", d)
	}
}

func TestLevenshtein_OneEdit(t *testing.T) {
	if d := levenshtein("kitten", "sitten"); d != 1 {
		t.Errorf("expected 1, got %d", d)
	}
}
