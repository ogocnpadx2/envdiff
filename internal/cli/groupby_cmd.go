package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

type groupByArgs struct {
	leftPath  string
	rightPath string
	prefix    string // optional: only show a specific prefix
}

func parseGroupByArgs(args []string) (groupByArgs, error) {
	fs := flag.NewFlagSet("groupby", flag.ContinueOnError)
	prefix := fs.String("prefix", "", "filter output to a specific prefix (e.g. DB_)")
	if err := fs.Parse(args); err != nil {
		return groupByArgs{}, err
	}
	if fs.NArg() < 2 {
		return groupByArgs{}, fmt.Errorf("usage: envdiff groupby <left> <right> [--prefix PREFIX]")
	}
	return groupByArgs{
		leftPath:  fs.Arg(0),
		rightPath: fs.Arg(1),
		prefix:    *prefix,
	}, nil
}

// RunGroupBy parses two .env files, diffs them, and prints keys grouped by prefix.
func RunGroupBy(args []string, out io.Writer) error {
	a, err := parseGroupByArgs(args)
	if err != nil {
		return err
	}

	left, err := parser.ParseFile(a.leftPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", a.leftPath, err)
	}
	right, err := parser.ParseFile(a.rightPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", a.rightPath, err)
	}

	result := diff.Compare(left, right)
	groups := diff.GroupResult(result)

	if len(groups) == 0 {
		fmt.Fprintln(out, "No differences found.")
		return nil
	}

	for _, g := range groups {
		if a.prefix != "" && !strings.EqualFold(g.Prefix, a.prefix) {
			continue
		}
		label := g.Prefix
		if label == "" {
			label = "(no prefix)"
		}
		fmt.Fprintf(out, "[%s]\n", label)
		for _, k := range g.Keys {
			fmt.Fprintf(out, "  %s\n", k)
		}
	}
	return nil
}

// RunGroupByMain is the entry-point wired into main.
func RunGroupByMain(args []string) {
	if err := RunGroupBy(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
