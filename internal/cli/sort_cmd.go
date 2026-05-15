package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

type sortArgs struct {
	path       string
	field      diff.SortField
	order      diff.SortOrder
	ignoreCase bool
}

func parseSortArgs(args []string) (sortArgs, error) {
	if len(args) < 1 {
		return sortArgs{}, fmt.Errorf("usage: envdiff sort <file> [--field key|value|length] [--order asc|desc] [--ignore-case]")
	}

	a := sortArgs{
		path:  args[0],
		field: diff.SortByKey,
		order: diff.SortAsc,
	}

	for i := 1; i < len(args); i++ {
		switch {
		case args[i] == "--ignore-case":
			a.ignoreCase = true
		case strings.HasPrefix(args[i], "--field="):
			v := strings.TrimPrefix(args[i], "--field=")
			f, err := diff.ParseSortField(v)
			if err != nil {
				return sortArgs{}, err
			}
			a.field = f
		case strings.HasPrefix(args[i], "--order="):
			v := strings.TrimPrefix(args[i], "--order=")
			o, err := diff.ParseSortOrder(v)
			if err != nil {
				return sortArgs{}, err
			}
			a.order = o
		default:
			return sortArgs{}, fmt.Errorf("unknown flag: %s", args[i])
		}
	}
	return a, nil
}

// RunSort executes the sort command writing output to w.
func RunSort(args []string, w io.Writer) error {
	a, err := parseSortArgs(args)
	if err != nil {
		return err
	}

	env, err := parser.ParseFile(a.path)
	if err != nil {
		return fmt.Errorf("parse %s: %w", a.path, err)
	}

	opts := diff.SortOptions{
		Field:      a.field,
		Order:      a.order,
		IgnoreCase: a.ignoreCase,
	}

	entries := diff.SortEnv(env, opts)
	for _, e := range entries {
		fmt.Fprintf(w, "%s=%s\n", e.Key, e.Value)
	}
	return nil
}

// RunSortMain is the entry point called from main dispatch.
func RunSortMain(args []string) {
	if err := RunSort(args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
