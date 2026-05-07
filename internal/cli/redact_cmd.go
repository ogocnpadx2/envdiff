package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

type redactArgs struct {
	leftPath    string
	rightPath   string
	patterns    string
	placeholder string
	format      string
}

func parseRedactArgs(args []string) (redactArgs, error) {
	fs := flag.NewFlagSet("redact", flag.ContinueOnError)
	patterns := fs.String("patterns", "", "comma-separated key patterns to redact (default: common secrets)")
	placeholder := fs.String("placeholder", diff.DefaultPlaceholder, "value substituted for redacted keys")
	format := fs.String("format", "text", "output format: text, json, markdown")

	if err := fs.Parse(args); err != nil {
		return redactArgs{}, err
	}
	positional := fs.Args()
	if len(positional) < 2 {
		return redactArgs{}, fmt.Errorf("usage: envdiff redact [flags] <left.env> <right.env>")
	}
	return redactArgs{
		leftPath:    positional[0],
		rightPath:   positional[1],
		patterns:    *patterns,
		placeholder: *placeholder,
		format:      *format,
	}, nil
}

// RunRedact compares two env files and prints the diff with sensitive values redacted.
func RunRedact(args []string, out io.Writer) error {
	ra, err := parseRedactArgs(args)
	if err != nil {
		return err
	}

	left, err := parser.ParseFile(ra.leftPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", ra.leftPath, err)
	}
	right, err := parser.ParseFile(ra.rightPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", ra.rightPath, err)
	}

	result := diff.Compare(left, right)
	opts := diff.ParseRedactOptions(ra.patterns, ra.placeholder)
	redacted := diff.RedactResult(result, opts)

	fmt_, err := diff.ParseFormat(ra.format)
	if err != nil {
		return err
	}
	return diff.PrintFormatted(out, redacted, ra.leftPath, ra.rightPath, fmt_)
}

// RunRedactMain is the entry point called from main for the redact subcommand.
func RunRedactMain(args []string) {
	if err := RunRedact(args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
