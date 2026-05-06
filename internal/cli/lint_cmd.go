package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

type lintArgs struct {
	path  string
	rules []diff.LintRule
}

func parseLintArgs(args []string) (lintArgs, error) {
	if len(args) < 1 {
		return lintArgs{}, fmt.Errorf("usage: envdiff lint <file> [--rules=<rule1,rule2>]")
	}

	path := args[0]
	rulesRaw := ""

	for _, arg := range args[1:] {
		if strings.HasPrefix(arg, "--rules=") {
			rulesRaw = strings.TrimPrefix(arg, "--rules=")
		}
	}

	rules, err := diff.ParseLintRules(rulesRaw)
	if err != nil {
		return lintArgs{}, fmt.Errorf("invalid rules: %w", err)
	}

	return lintArgs{path: path, rules: rules}, nil
}

// RunLint executes the lint subcommand.
func RunLint(args []string, stdout, stderr io.Writer) int {
	la, err := parseLintArgs(args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	env, err := parser.ParseFile(la.path)
	if err != nil {
		fmt.Fprintf(stderr, "error reading %s: %v\n", la.path, err)
		return 1
	}

	result := diff.Lint(env, la.rules)

	if result.Clean() {
		fmt.Fprintf(stdout, "OK: no lint violations in %s\n", la.path)
		return 0
	}

	fmt.Fprintf(stdout, "Lint violations in %s:\n", la.path)
	for _, v := range result.Violations {
		fmt.Fprintf(stdout, "  %s\n", v.String())
	}
	return 1
}

// RunLintMain is the entry point called from main when subcommand is "lint".
func RunLintMain(args []string) int {
	return RunLint(args, os.Stdout, os.Stderr)
}
