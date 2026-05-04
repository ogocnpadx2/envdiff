package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents the key-value pairs parsed from a .env file.
type EnvMap map[string]string

// ParseFile reads a .env file and returns its key-value pairs.
// It skips blank lines and comments (lines starting with '#').
// It returns an error if the file cannot be opened or contains malformed lines.
func ParseFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parser: cannot open file %q: %w", path, err)
	}
	defer f.Close()

	env := make(EnvMap)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("parser: %s line %d: %w", path, lineNum, err)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parser: error reading %q: %w", path, err)
	}

	return env, nil
}

// parseLine splits a single "KEY=VALUE" line into its components.
func parseLine(line string) (string, string, error) {
	idx := strings.IndexByte(line, '=')
	if idx < 0 {
		return "", "", fmt.Errorf("invalid line %q: missing '='")
	}

	key := strings.TrimSpace(line[:idx])
	if key == "" {
		return "", "", fmt.Errorf("invalid line %q: empty key", line)
	}

	value := strings.TrimSpace(line[idx+1:])
	value = stripQuotes(value)

	return key, value, nil
}

// stripQuotes removes surrounding single or double quotes from a value.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
