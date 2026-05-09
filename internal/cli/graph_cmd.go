package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

type graphArgs struct {
	path      string
	showCycles bool
}

func parseGraphArgs(args []string) (graphArgs, error) {
	if len(args) == 0 {
		return graphArgs{}, fmt.Errorf("usage: envdiff graph <file> [--cycles]")
	}
	ga := graphArgs{path: args[0]}
	for _, a := range args[1:] {
		if a == "--cycles" {
			ga.showCycles = true
		} else {
			return graphArgs{}, fmt.Errorf("unknown flag: %s", a)
		}
	}
	return ga, nil
}

// RunGraph prints the dependency graph for a single .env file.
func RunGraph(args []string, out io.Writer) error {
	ga, err := parseGraphArgs(args)
	if err != nil {
		return err
	}

	env, err := parser.ParseFile(ga.path)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	g := diff.BuildDependencyGraph(env)

	if ga.showCycles {
		cycles := g.CyclicKeys()
		if len(cycles) == 0 {
			fmt.Fprintln(out, "No cycles detected.")
		} else {
			fmt.Fprintln(out, "Cyclic keys:")
			for _, k := range cycles {
				fmt.Fprintf(out, "  %s\n", k)
			}
		}
		return nil
	}

	entries := g.SortedEntries()
	hasDeps := false
	for _, e := range entries {
		if len(e.Deps) > 0 {
			hasDeps = true
			fmt.Fprintf(out, "%s -> %s\n", e.Key, strings.Join(e.Deps, ", "))
		}
	}
	if !hasDeps {
		fmt.Fprintln(out, "No inter-key dependencies found.")
	}
	return nil
}

// RunGraphMain is the top-level entry point called from main dispatch.
func RunGraphMain(args []string) {
	if err := RunGraph(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
