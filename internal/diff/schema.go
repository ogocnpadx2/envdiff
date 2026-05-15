package diff

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Schema represents a required set of keys with optional descriptions.
type Schema struct {
	Keys map[string]string // key -> description (may be empty)
}

// LoadSchema reads a schema file where each line is:
//
//	KEY
//	KEY=description of what this key is for
//
// Blank lines and lines starting with '#' are ignored.
func LoadSchema(path string) (*Schema, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open schema %q: %w", path, err)
	}
	defer f.Close()

	s := &Schema{Keys: make(map[string]string)}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		if key == "" {
			continue
		}
		desc := ""
		if len(parts) == 2 {
			desc = strings.TrimSpace(parts[1])
		}
		s.Keys[key] = desc
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read schema %q: %w", path, err)
	}
	return s, nil
}

// SchemaViolation describes a key required by the schema but absent from an env map.
type SchemaViolation struct {
	Key         string
	Description string
}

// ValidateAgainstSchema returns violations for any schema-required keys missing from env.
func ValidateAgainstSchema(schema *Schema, env map[string]string) []SchemaViolation {
	var violations []SchemaViolation
	for key, desc := range schema.Keys {
		if _, ok := env[key]; !ok {
			violations = append(violations, SchemaViolation{Key: key, Description: desc})
		}
	}
	sortViolations(violations)
	return violations
}

// ExtraKeys returns keys present in env that are not defined in the schema.
// This is useful for detecting undocumented or unexpected environment variables.
func ExtraKeys(schema *Schema, env map[string]string) []string {
	var extras []string
	for key := range env {
		if _, ok := schema.Keys[key]; !ok {
			extras = append(extras, key)
		}
	}
	sortStrings(extras)
	return extras
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}

func sortViolations(v []SchemaViolation) {
	for i := 1; i < len(v); i++ {
		for j := i; j > 0 && v[j].Key < v[j-1].Key; j-- {
			v[j], v[j-1] = v[j-1], v[j]
		}
	}
}
