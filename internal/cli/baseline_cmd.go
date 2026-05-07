package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

type baselineArgs struct {
	leftPath  string
	rightPath string
	baseline  string
	save      bool
}

func parseBaselineArgs(args []string) (baselineArgs, error) {
	if len(args) < 2 {
		return baselineArgs{}, fmt.Errorf("usage: envdiff baseline <left> <right> --baseline <file> [--save]")
	}
	a := baselineArgs{
		leftPath:  args[0],
		rightPath: args[1],
		baseline:  ".envdiff-baseline.json",
	}
	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "--save":
			a.save = true
		case "--baseline":
			if i+1 >= len(args) {
				return baselineArgs{}, fmt.Errorf("--baseline requires a file path")
			}
			i++
			a.baseline = args[i]
		}
	}
	return a, nil
}

// RunBaseline saves or compares against a baseline diff.
func RunBaseline(args []string, stdout, stderr io.Writer) int {
	a, err := parseBaselineArgs(args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	left, err := parser.ParseFile(a.leftPath)
	if err != nil {
		fmt.Fprintf(stderr, "error reading %s: %v\n", a.leftPath, err)
		return 1
	}
	right, err := parser.ParseFile(a.rightPath)
	if err != nil {
		fmt.Fprintf(stderr, "error reading %s: %v\n", a.rightPath, err)
		return 1
	}
	result := diff.Compare(left, right)

	if a.save {
		if err := diff.SaveBaseline(a.baseline, a.leftPath, a.rightPath, result); err != nil {
			fmt.Fprintf(stderr, "error saving baseline: %v\n", err)
			return 1
		}
		fmt.Fprintf(stdout, "Baseline saved to %s\n", a.baseline)
		return 0
	}

	baseline, err := diff.LoadBaseline(a.baseline)
	if err != nil {
		fmt.Fprintf(stderr, "error loading baseline: %v\n", err)
		return 1
	}

	newIssues, resolved := diff.DiffAgainstBaseline(baseline, result)

	if len(resolved.MissingInLeft)+len(resolved.MissingInRight)+len(resolved.Mismatched) > 0 {
		fmt.Fprintln(stdout, "Resolved since baseline:")
		diff.PrintReport(stdout, resolved)
	}
	if len(newIssues.MissingInLeft)+len(newIssues.MissingInRight)+len(newIssues.Mismatched) > 0 {
		fmt.Fprintln(stdout, "New issues since baseline:")
		diff.PrintReport(stdout, newIssues)
		return 1
	}
	fmt.Fprintln(stdout, "No new issues since baseline.")
	return 0
}

// RunBaselineMain is the entry point called from main dispatch.
func RunBaselineMain(args []string) int {
	return RunBaseline(args, os.Stdout, os.Stderr)
}
