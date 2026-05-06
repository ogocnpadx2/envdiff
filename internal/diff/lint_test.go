package diff

import (
	"testing"
)

func TestParseLintRules_Default(t *testing.T) {
	rules, err := ParseLintRules("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) == 0 {
		t.Error("expected default rules, got none")
	}
}

func TestParseLintRules_Valid(t *testing.T) {
	rules, err := ParseLintRules("uppercase,non-empty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(rules))
	}
}

func TestParseLintRules_Invalid(t *testing.T) {
	_, err := ParseLintRules("uppercase,bad-rule")
	if err == nil {
		t.Error("expected error for unknown rule")
	}
}

func TestLint_Clean(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	rules, _ := ParseLintRules("uppercase,no-spaces,non-empty")
	result := Lint(env, rules)
	if !result.Clean() {
		t.Errorf("expected clean result, got violations: %v", result.Violations)
	}
}

func TestLint_UppercaseViolation(t *testing.T) {
	env := map[string]string{"app_name": "myapp"}
	rules := []LintRule{LintRuleUppercase}
	result := Lint(env, rules)
	if result.Clean() {
		t.Error("expected uppercase violation")
	}
	if result.Violations[0].Rule != LintRuleUppercase {
		t.Errorf("expected uppercase rule, got %s", result.Violations[0].Rule)
	}
}

func TestLint_NoSpacesViolation(t *testing.T) {
	env := map[string]string{"KEY": "  value  "}
	rules := []LintRule{LintRuleNoSpaces}
	result := Lint(env, rules)
	if result.Clean() {
		t.Error("expected no-spaces violation")
	}
}

func TestLint_NonEmptyViolation(t *testing.T) {
	env := map[string]string{"KEY": ""}
	rules := []LintRule{LintRuleNonEmpty}
	result := Lint(env, rules)
	if result.Clean() {
		t.Error("expected non-empty violation")
	}
}

func TestLint_NoQuotesViolation(t *testing.T) {
	env := map[string]string{"KEY": `"value"`}
	rules := []LintRule{LintRuleNoQuotes}
	result := Lint(env, rules)
	if result.Clean() {
		t.Error("expected no-quotes violation")
	}
}

func TestLint_MultipleViolations(t *testing.T) {
	env := map[string]string{
		"lower_key": "",
		"GOOD_KEY":  "ok",
	}
	rules := []LintRule{LintRuleUppercase, LintRuleNonEmpty}
	result := Lint(env, rules)
	// lower_key violates both uppercase and non-empty
	if len(result.Violations) < 2 {
		t.Errorf("expected at least 2 violations, got %d", len(result.Violations))
	}
}
