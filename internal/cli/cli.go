package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// Run is the entry point for the CLI. It parses args, runs the diff, and
// writes output to w. It returns a non-zero exit code when differences exist.
func Run(args []string, w io.Writer) int {
	flags, ok := parseArgs(args, w)
	if !ok {
		return 2
	}

	left, err := parser.ParseFile(flags.leftFile)
	if err != nil {
		fmt.Fprintf(w, "error reading %s: %v\n", flags.leftFile, err)
		return 2
	}
	right, err := parser.ParseFile(flags.rightFile)
	if err != nil {
		fmt.Fprintf(w, "error reading %s: %v\n", flags.rightFile, err)
		return 2
	}

	result := diff.Compare(left, right)

	if flags.quiet {
		if result.IsClean() {
			return 0
		}
		return 1
	}

	fmt := flags.format
	leftName := filepath.Base(flags.leftFile)
	rightName := filepath.Base(flags.rightFile)

	if err := diff.PrintFormatted(w, result, leftName, rightName, fmt); err != nil {
		fmt2 := fmt // shadow to avoid collision
		_ = fmt2
		fmt2 = diff.FormatText
		_ = fmt2
		// fallback: write error
		io.WriteString(w, "error formatting output: "+err.Error()+"\n")
		return 2
	}

	if !result.IsClean() {
		return 1
	}
	return 0
}

type cliFlags struct {
	leftFile  string
	rightFile string
	quiet     bool
	format    diff.OutputFormat
}

func parseArgs(args []string, w io.Writer) (cliFlags, bool) {
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)
	fs.SetOutput(w)

	quiet := fs.Bool("quiet", false, "exit with code only, no output")
	formatStr := fs.String("format", "text", "output format: text, json, markdown")

	if err := fs.Parse(args); err != nil {
		return cliFlags{}, false
	}

	if fs.NArg() < 2 {
		fmt.Fprintln(w, "usage: envdiff [--quiet] [--format=text|json|markdown] <left.env> <right.env>")
		return cliFlags{}, false
	}

	fmt2, err := diff.ParseFormat(*formatStr)
	if err != nil {
		fmt.Fprintln(w, err)
		return cliFlags{}, false
	}

	_ = os.Stderr // suppress unused import

	return cliFlags{
		leftFile:  fs.Arg(0),
		rightFile: fs.Arg(1),
		quiet:     *quiet,
		format:    fmt2,
	}, true
}
