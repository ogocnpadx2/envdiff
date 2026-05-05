package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// Config holds parsed CLI options.
type Config struct {
	LeftFile  string
	RightFile string
	Quiet     bool
	NoColor   bool
}

// Run is the entrypoint for the CLI. It parses args, runs the diff, and writes
// output to out/errOut. Returns a non-nil error when the program should exit
// with a non-zero status.
func Run(args []string, out, errOut io.Writer) error {
	cfg, err := parseArgs(args, errOut)
	if err != nil {
		return err
	}

	left, err := parser.ParseFile(cfg.LeftFile)
	if err != nil {
		return fmt.Errorf("reading %s: %w", cfg.LeftFile, err)
	}

	right, err := parser.ParseFile(cfg.RightFile)
	if err != nil {
		return fmt.Errorf("reading %s: %w", cfg.RightFile, err)
	}

	report := diff.Compare(left, right)

	if cfg.Quiet {
		if !report.Clean() {
			return errors.New("environments differ")
		}
		return nil
	}

	diff.PrintReport(out, cfg.LeftFile, cfg.RightFile, report)

	if !report.Clean() {
		return errors.New("environments differ")
	}
	return nil
}

func parseArgs(args []string, errOut io.Writer) (*Config, error) {
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)
	fs.SetOutput(errOut)

	quiet := fs.Bool("quiet", false, "suppress output; exit 1 if differences found")
	noColor := fs.Bool("no-color", false, "disable colored output")

	fs.Usage = func() {
		fmt.Fprintf(errOut, "Usage: envdiff [options] <file1> <file2>\n\nOptions:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if fs.NArg() != 2 {
		fs.Usage()
		return nil, errors.New("exactly two .env files are required")
	}

	return &Config{
		LeftFile:  fs.Arg(0),
		RightFile: fs.Arg(1),
		Quiet:     *quiet,
		NoColor:   *noColor,
	}, nil
}
