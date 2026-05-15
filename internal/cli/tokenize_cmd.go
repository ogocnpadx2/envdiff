package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

type tokenizeArgs struct {
	leftPath  string
	rightPath string
	delimiter string
	lowercase bool
}

func parseTokenizeArgs(args []string) (tokenizeArgs, error) {
	if len(args) < 2 {
		return tokenizeArgs{}, fmt.Errorf("usage: envdiff tokenize <left> <right> [--delimiter=,] [--no-lowercase]")
	}
	a := tokenizeArgs{
		leftPath:  args[0],
		rightPath: args[1],
		delimiter: ",",
		lowercase: true,
	}
	for _, flag := range args[2:] {
		if strings.HasPrefix(flag, "--delimiter=") {
			a.delimiter = strings.TrimPrefix(flag, "--delimiter=")
		} else if flag == "--no-lowercase" {
			a.lowercase = false
		} else {
			return tokenizeArgs{}, fmt.Errorf("unknown flag: %s", flag)
		}
	}
	return a, nil
}

// RunTokenize is the entry point for the tokenize subcommand.
func RunTokenize(args []string, out io.Writer) error {
	a, err := parseTokenizeArgs(args)
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

	opts := diff.TokenizeOptions{
		Delimiter: a.delimiter,
		Lowercase: a.lowercase,
	}
	results := diff.TokenizeValues(left, right, opts)

	if len(results) == 0 {
		fmt.Fprintln(out, "all token values match")
		return nil
	}

	for _, r := range results {
		fmt.Fprintf(out, "[%s]\n", r.Key)
		if len(r.OnlyInLeft) > 0 {
			fmt.Fprintf(out, "  only in left:  %s\n", strings.Join(r.OnlyInLeft, ", "))
		}
		if len(r.OnlyInRight) > 0 {
			fmt.Fprintf(out, "  only in right: %s\n", strings.Join(r.OnlyInRight, ", "))
		}
		if len(r.Shared) > 0 {
			fmt.Fprintf(out, "  shared:        %s\n", strings.Join(r.Shared, ", "))
		}
	}
	return nil
}

// RunTokenizeMain is called from main for the tokenize subcommand.
func RunTokenizeMain(args []string) {
	if err := RunTokenize(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
