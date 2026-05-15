package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

type pivotArgs struct {
	files  []string // "label=path" pairs
}

func parsePivotArgs(args []string) (pivotArgs, error) {
	if len(args) < 2 {
		return pivotArgs{}, fmt.Errorf("pivot requires at least two label=path arguments")
	}
	for _, a := range args {
		if !strings.Contains(a, "=") {
			return pivotArgs{}, fmt.Errorf("argument %q must be in label=path format", a)
		}
	}
	return pivotArgs{files: args}, nil
}

// RunPivot builds and prints a pivot table from multiple labelled env files.
func RunPivot(args []string, out io.Writer) error {
	pa, err := parsePivotArgs(args)
	if err != nil {
		return err
	}

	envs := make(map[string]map[string]string, len(pa.files))
	for _, arg := range pa.files {
		parts := strings.SplitN(arg, "=", 2)
		label, path := parts[0], parts[1]
		kv, err := parser.ParseFile(path)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}
		envs[label] = kv
	}

	pt := diff.BuildPivot(envs)
	fmt.Fprint(out, diff.FormatPivotText(pt))
	return nil
}

// RunPivotMain is the entry point called from main dispatch.
func RunPivotMain(args []string) {
	if err := RunPivot(args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
