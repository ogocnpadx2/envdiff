package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// Run is the entry point for the CLI. It parses args, runs the comparison, and
// writes output to out. Returns an exit code.
func Run(args []string, out io.Writer) int {
	opts, err := parseArgs(args, out)
	if err != nil {
		return 2
	}

	left, err := parser.ParseFile(opts.leftFile)
	if err != nil {
		fmt.Fprintf(out, "error reading %s: %v\n", opts.leftFile, err)
		return 1
	}
	right, err := parser.ParseFile(opts.rightFile)
	if err != nil {
		fmt.Fprintf(out, "error reading %s: %v\n", opts.rightFile, err)
		return 1
	}

	result := diff.Compare(left, right)
	result = diff.FilterResult(result, diff.FilterOptions{
		OnlyMissing:    opts.onlyMissing,
		OnlyMismatched: opts.onlyMismatched,
		KeyPrefix:      opts.keyPrefix,
	})

	// Schema validation
	if opts.schemaFile != "" {
		schema, err := diff.LoadSchema(opts.schemaFile)
		if err != nil {
			fmt.Fprintf(out, "error loading schema %s: %v\n", opts.schemaFile, err)
			return 1
		}
		violations := diff.ValidateAgainstSchema(schema, left)
		if len(violations) > 0 {
			fmt.Fprintf(out, "Schema violations in %s:\n", opts.leftFile)
			for _, v := range violations {
				if v.Description != "" {
					fmt.Fprintf(out, "  missing: %s (%s)\n", v.Key, v.Description)
				} else {
					fmt.Fprintf(out, "  missing: %s\n", v.Key)
				}
			}
		}
	}

	format, err := diff.ParseFormat(opts.format)
	if err != nil {
		fmt.Fprintf(out, "error: %v\n", err)
		return 2
	}

	if opts.quiet {
		severity := diff.MaxSeverity(diff.ClassifyResult(result))
		if severity == diff.SeverityNone {
			return 0
		}
		return 1
	}

	if err := diff.PrintFormatted(out, result, opts.leftFile, opts.rightFile, format); err != nil {
		fmt.Fprintf(out, "error: %v\n", err)
		return 1
	}

	if diff.MaxSeverity(diff.ClassifyResult(result)) > diff.SeverityNone {
		return 1
	}
	return 0
}

type options struct {
	leftFile       string
	rightFile      string
	quiet          bool
	onlyMissing    bool
	onlyMismatched bool
	keyPrefix      string
	format         string
	schemaFile     string
}

func parseArgs(args []string, out io.Writer) (*options, error) {
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)
	fs.SetOutput(out)

	opts := &options{}
	fs.BoolVar(&opts.quiet, "quiet", false, "exit 1 if differences found, no output")
	fs.BoolVar(&opts.onlyMissing, "only-missing", false, "show only missing keys")
	fs.BoolVar(&opts.onlyMismatched, "only-mismatched", false, "show only mismatched values")
	fs.StringVar(&opts.keyPrefix, "prefix", "", "filter keys by prefix")
	fs.StringVar(&opts.format, "format", "text", "output format: text, json, markdown")
	fs.StringVar(&opts.schemaFile, "schema", "", "path to schema file of required keys")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if fs.NArg() < 2 {
		fmt.Fprintf(out, "usage: envdiff [flags] <file1> <file2>\n")
		fs.PrintDefaults()
		return nil, fmt.Errorf("missing arguments")
	}
	opts.leftFile = fs.Arg(0)
	opts.rightFile = fs.Arg(1)
	_ = os.DevNull
	return opts, nil
}
