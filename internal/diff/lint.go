package diff

import (
	"fmt"
	"sort"
	"strings"
)

// LintRule represents a validation rule applied to env keys/values.
type LintRule string

const (
	LintRuleUppercase  LintRule = "uppercase"   // keys must be uppercase
	LintRuleNoSpaces   LintRule = "no-spaces"   // values must not contain leading/trailing spaces
	LintRuleNonEmpty   LintRule = "non-empty"   // values must not be empty
	LintRuleNoQuotes   LintRule = "no-quotes"   // values must not be wrapped in quotes
)

// LintViolation describes a single linting issue.
type LintViolation struct {
	Key     string
	Rule    LintRule
	Message string
}

func (v LintViolation) String() string {
	return fmt.Sprintf("[%s] %s: %s", v.Rule, v.Key, v.Message)
}

// LintResult holds all violations found during linting.
type LintResult struct {
	Violations []LintViolation
}

func (r LintResult) Clean() bool {
	return len(r.Violations) == 0
}

// ParseLintRules parses a comma-separated list of rule names.
func ParseLintRules(raw string) ([]LintRule, error) {
	if raw == "" {
		return defaultLintRules(), nil
	}
	parts := strings.Split(raw, ",")
	rules := make([]LintRule, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		switch LintRule(p) {
		case LintRuleUppercase, LintRuleNoSpaces, LintRuleNonEmpty, LintRuleNoQuotes:
			rules = append(rules, LintRule(p))
		default:
			return nil, fmt.Errorf("unknown lint rule: %q", p)
		}
	}
	return rules, nil
}

func defaultLintRules() []LintRule {
	return []LintRule{LintRuleUppercase, LintRuleNoSpaces, LintRuleNonEmpty}
}

// Lint checks the given env map against the provided rules.
func Lint(env map[string]string, rules []LintRule) LintResult {
	var violations []LintViolation
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]
		for _, rule := range rules {
			switch rule {
			case LintRuleUppercase:
				if k != strings.ToUpper(k) {
					violations = append(violations, LintViolation{Key: k, Rule: rule, Message: "key should be uppercase"})
				}
			case LintRuleNoSpaces:
				if strings.TrimSpace(v) != v {
					violations = append(violations, LintViolation{Key: k, Rule: rule, Message: "value has leading or trailing spaces"})
				}
			case LintRuleNonEmpty:
				if v == "" {
					violations = append(violations, LintViolation{Key: k, Rule: rule, Message: "value is empty"})
				}
			case LintRuleNoQuotes:
				if (strings.HasPrefix(v, `"`) && strings.HasSuffix(v, `"`)) ||
					(strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'")) {
					violations = append(violations, LintViolation{Key: k, Rule: rule, Message: "value is wrapped in quotes"})
				}
			}
		}
	}
	return LintResult{Violations: violations}
}
