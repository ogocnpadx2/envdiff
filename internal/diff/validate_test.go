package diff

import (
	"testing"
)

func TestParseValueRules_Default(t *testing.T) {
	rules, err := ParseValueRules(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) == 0 {
		t.Fatal("expected at least one default rule")
	}
}

func TestParseValueRules_Valid(t *testing.T) {
	rules, err := ParseValueRules([]string{"no-empty"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 || rules[0].Name != "no-empty" {
		t.Fatalf("unexpected rules: %+v", rules)
	}
}

func TestParseValueRules_Invalid(t *testing.T) {
	_, err := ParseValueRules([]string{"not-a-rule"})
	if err == nil {
		t.Fatal("expected error for unknown rule")
	}
}

func TestValidateValues_Clean(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	rules, _ := ParseValueRules(nil)
	result := ValidateValues(env, rules)
	if !result.IsClean() {
		t.Fatalf("expected clean result, got violations: %+v", result.Violations)
	}
	if result.Summary() != "all values passed validation" {
		t.Errorf("unexpected summary: %s", result.Summary())
	}
}

func TestValidateValues_EmptyValue(t *testing.T) {
	env := map[string]string{
		"HOST": "",
		"PORT": "5432",
	}
	rules, _ := ParseValueRules(nil)
	result := ValidateValues(env, rules)
	if result.IsClean() {
		t.Fatal("expected violations for empty value")
	}
	if result.Violations[0].Key != "HOST" {
		t.Errorf("expected HOST violation, got %s", result.Violations[0].Key)
	}
	if result.Violations[0].Rule != "no-empty" {
		t.Errorf("expected no-empty rule, got %s", result.Violations[0].Rule)
	}
}

func TestValidateValues_WhitespaceOnly(t *testing.T) {
	env := map[string]string{
		"SECRET": "   ",
	}
	rules, _ := ParseValueRules([]string{"no-whitespace-only"})
	result := ValidateValues(env, rules)
	if result.IsClean() {
		t.Fatal("expected violation for whitespace-only value")
	}
	if result.Violations[0].Rule != "no-whitespace-only" {
		t.Errorf("unexpected rule: %s", result.Violations[0].Rule)
	}
}

func TestValidateValues_SortedOutput(t *testing.T) {
	env := map[string]string{
		"Z_KEY": "",
		"A_KEY": "",
		"M_KEY": "",
	}
	rules, _ := ParseValueRules([]string{"no-empty"})
	result := ValidateValues(env, rules)
	if len(result.Violations) != 3 {
		t.Fatalf("expected 3 violations, got %d", len(result.Violations))
	}
	if result.Violations[0].Key != "A_KEY" || result.Violations[1].Key != "M_KEY" || result.Violations[2].Key != "Z_KEY" {
		t.Errorf("violations not sorted: %+v", result.Violations)
	}
}

func TestValidationResult_Summary_WithViolations(t *testing.T) {
	r := ValidationResult{
		Violations: []ValidationViolation{{Key: "K", Rule: "no-empty"}},
	}
	if r.Summary() != "1 validation violation(s) found" {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}
