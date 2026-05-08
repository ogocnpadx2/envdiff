package diff

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// ValueRule defines a validation rule applied to env values.
type ValueRule struct {
	Name    string
	Pattern *regexp.Regexp
	Message string
}

// ValidationViolation describes a single failed validation.
type ValidationViolation struct {
	Key     string
	Value   string
	Rule    string
	Message string
}

// ValidationResult holds all violations found during validation.
type ValidationResult struct {
	Violations []ValidationViolation
}

// IsClean returns true when no violations were found.
func (r ValidationResult) IsClean() bool {
	return len(r.Violations) == 0
}

// Summary returns a short human-readable summary.
func (r ValidationResult) Summary() string {
	if r.IsClean() {
		return "all values passed validation"
	}
	return fmt.Sprintf("%d validation violation(s) found", len(r.Violations))
}

var builtinRules = []ValueRule{
	{
		Name:    "no-empty",
		Pattern: regexp.MustCompile(`^.+$`),
		Message: "value must not be empty",
	},
	{
		Name:    "no-whitespace-only",
		Pattern: regexp.MustCompile(`\S`),
		Message: "value must not be whitespace-only",
	},
}

// ParseValueRules returns a slice of ValueRule by name, or an error for
// unknown rule names. Passing nil or empty slice returns the builtin rules.
func ParseValueRules(names []string) ([]ValueRule, error) {
	if len(names) == 0 {
		return builtinRules, nil
	}
	index := make(map[string]ValueRule, len(builtinRules))
	for _, r := range builtinRules {
		index[r.Name] = r
	}
	var out []ValueRule
	for _, n := range names {
		r, ok := index[strings.TrimSpace(n)]
		if !ok {
			return nil, fmt.Errorf("unknown validation rule: %q", n)
		}
		out = append(out, r)
	}
	return out, nil
}

// ValidateValues checks every key/value pair in env against the provided rules
// and returns a ValidationResult.
func ValidateValues(env map[string]string, rules []ValueRule) ValidationResult {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var violations []ValidationViolation
	for _, k := range keys {
		v := env[k]
		for _, rule := range rules {
			if !rule.Pattern.MatchString(v) {
				violations = append(violations, ValidationViolation{
					Key:     k,
					Value:   v,
					Rule:    rule.Name,
					Message: rule.Message,
				})
			}
		}
	}
	return ValidationResult{Violations: violations}
}
