package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

type snapshotArgs struct {
	leftFile  string
	rightFile string
	outputPath string
	comparePath string
}

func parseSnapshotArgs(args []string) (snapshotArgs, error) {
	var a snapshotArgs
	rest := args
	for len(rest) > 0 {
		switch rest[0] {
		case "--save":
			if len(rest) < 2 {
				return a, fmt.Errorf("--save requires a path argument")
			}
			a.outputPath = rest[1]
			rest = rest[2:]
		case "--compare":
			if len(rest) < 2 {
				return a, fmt.Errorf("--compare requires a path argument")
			}
			a.comparePath = rest[1]
			rest = rest[2:]
		default:
			if a.leftFile == "" {
				a.leftFile = rest[0]
			} else if a.rightFile == "" {
				a.rightFile = rest[0]
			} else {
				return a, fmt.Errorf("unexpected argument: %s", rest[0])
			}
			rest = rest[1:]
		}
	}
	if a.leftFile == "" || a.rightFile == "" {
		return a, fmt.Errorf("snapshot requires two .env file paths")
	}
	if a.outputPath == "" && a.comparePath == "" {
		return a, fmt.Errorf("snapshot requires --save <path> or --compare <path>")
	}
	return a, nil
}

// RunSnapshot handles the snapshot subcommand.
func RunSnapshot(args []string, stdout, stderr io.Writer) int {
	a, err := parseSnapshotArgs(args)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	left, err := parser.ParseFile(a.leftFile)
	if err != nil {
		fmt.Fprintf(stderr, "error reading %s: %v\n", a.leftFile, err)
		return 1
	}
	right, err := parser.ParseFile(a.rightFile)
	if err != nil {
		fmt.Fprintf(stderr, "error reading %s: %v\n", a.rightFile, err)
		return 1
	}

	result := diff.Compare(left, right)

	if a.outputPath != "" {
		if err := diff.SaveSnapshot(a.outputPath, a.leftFile, a.rightFile, result); err != nil {
			fmt.Fprintf(stderr, "error saving snapshot: %v\n", err)
			return 1
		}
		fmt.Fprintf(stdout, "snapshot saved to %s\n", a.outputPath)
	}

	if a.comparePath != "" {
		before, err := diff.LoadSnapshot(a.comparePath)
		if err != nil {
			fmt.Fprintf(stderr, "error loading snapshot: %v\n", err)
			return 1
		}
		after := &diff.Snapshot{Result: result}
		delta := diff.DiffSnapshots(before, after)
		if len(delta.NewIssues) == 0 && len(delta.ResolvedIssues) == 0 {
			fmt.Fprintln(stdout, "no changes since snapshot")
			return 0
		}
		if len(delta.NewIssues) > 0 {
			fmt.Fprintf(stdout, "new issues: %v\n", delta.NewIssues)
		}
		if len(delta.ResolvedIssues) > 0 {
			fmt.Fprintf(stdout, "resolved: %v\n", delta.ResolvedIssues)
		}
		if len(delta.NewIssues) > 0 {
			return 1
		}
	}

	_ = os.Stdout // ensure os import used
	return 0
}
