package cli

import (
	"fmt"
	"io"
	"os"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

type profileArgs struct {
	path string
}

func parseProfileArgs(args []string) (profileArgs, error) {
	if len(args) < 1 {
		return profileArgs{}, fmt.Errorf("usage: envdiff profile <file>")
	}
	return profileArgs{path: args[0]}, nil
}

// RunProfile parses an env file and prints a type/shape profile of its values.
func RunProfile(args []string, out io.Writer) error {
	pa, err := parseProfileArgs(args)
	if err != nil {
		return err
	}

	env, err := parser.ParseFile(pa.path)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	result := diff.ProfileEnv(env)

	if len(result.Entries) == 0 {
		fmt.Fprintln(out, "No keys found.")
		return nil
	}

	fmt.Fprintf(out, "%-40s %-8s %s\n", "KEY", "TYPE", "LENGTH")
	fmt.Fprintf(out, "%-40s %-8s %s\n", "---", "----", "------")
	for _, e := range result.Entries {
		empty := ""
		if !e.NonEmpty {
			empty = " (empty)"
		}
		fmt.Fprintf(out, "%-40s %-8s %d%s\n", e.Key, e.Type, e.Length, empty)
	}
	return nil
}

// RunProfileMain is the entry point wired into the top-level CLI dispatcher.
func RunProfileMain(args []string) {
	if err := RunProfile(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
