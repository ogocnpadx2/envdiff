package cli

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

type suggestArgs struct {
	leftPath    string
	rightPath   string
	maxDistance int
}

func parseSuggestArgs(args []string) (suggestArgs, error) {
	if len(args) < 2 {
		return suggestArgs{}, fmt.Errorf("usage: envdiff suggest <left> <right> [--max-distance=N]")
	}
	sa := suggestArgs{
		leftPath:    args[0],
		rightPath:   args[1],
		maxDistance: 3,
	}
	for _, arg := range args[2:] {
		if v, ok := stripFlag(arg, "--max-distance"); ok {
			n, err := strconv.Atoi(v)
			if err != nil || n < 1 {
				return suggestArgs{}, fmt.Errorf("--max-distance must be a positive integer")
			}
			sa.maxDistance = n
		}
	}
	return sa, nil
}

// RunSuggest runs the suggest sub-command writing output to w.
func RunSuggest(args []string, w io.Writer) error {
	sa, err := parseSuggestArgs(args)
	if err != nil {
		return err
	}

	left, err := parser.ParseFile(sa.leftPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", sa.leftPath, err)
	}
	right, err := parser.ParseFile(sa.rightPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", sa.rightPath, err)
	}

	result := diff.Compare(left, right)
	sr := diff.Suggest(result, sa.maxDistance)

	if len(sr.Suggestions) == 0 {
		fmt.Fprintln(w, "No suggestions — no likely typos or renames detected.")
		return nil
	}

	fmt.Fprintf(w, "Suggestions (%d):\n", len(sr.Suggestions))
	for _, s := range sr.Suggestions {
		fmt.Fprintf(w, "  %-30s  ->  %s  (distance: %d)\n", s.MissingKey, s.SuggestedKey, s.Distance)
	}
	return nil
}

// RunSuggestMain is the entry-point wired from main.
func RunSuggestMain(args []string) {
	if err := RunSuggest(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// stripFlag parses "--key=value" or "--key value" style (key=value form only here).
func stripFlag(arg, name string) (string, bool) {
	prefix := name + "="
	if len(arg) > len(prefix) && arg[:len(prefix)] == prefix {
		return arg[len(prefix):], true
	}
	return "", false
}
