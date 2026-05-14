package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

type ignoreArgs struct {
	leftPath   string
	rightPath  string
	ignoreFile string
	ignoreKeys []string
}

func parseIgnoreArgs(args []string) (ignoreArgs, error) {
	fs := flag.NewFlagSet("ignore", flag.ContinueOnError)
	ignoreFile := fs.String("ignore-file", "", "path to .envignore file")
	ignoreKeys := fs.String("ignore-keys", "", "comma-separated list of keys to ignore")

	if err := fs.Parse(args); err != nil {
		return ignoreArgs{}, err
	}
	positional := fs.Args()
	if len(positional) < 2 {
		return ignoreArgs{}, fmt.Errorf("usage: envdiff ignore [flags] <left> <right>")
	}

	var keys []string
	if *ignoreKeys != "" {
		for _, k := range strings.Split(*ignoreKeys, ",") {
			if k = strings.TrimSpace(k); k != "" {
				keys = append(keys, k)
			}
		}
	}

	return ignoreArgs{
		leftPath:   positional[0],
		rightPath:  positional[1],
		ignoreFile: *ignoreFile,
		ignoreKeys: keys,
	}, nil
}

// RunIgnore executes the ignore-aware diff sub-command.
func RunIgnore(args []string, out io.Writer) error {
	ia, err := parseIgnoreArgs(args)
	if err != nil {
		return err
	}

	left, err := parser.ParseFile(ia.leftPath)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", ia.leftPath, err)
	}
	right, err := parser.ParseFile(ia.rightPath)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", ia.rightPath, err)
	}

	result := diff.Compare(left, right)

	opts := diff.IgnoreOptions{Keys: ia.ignoreKeys}
	if ia.ignoreFile != "" {
		fileOpts, err := diff.ParseIgnoreFile(ia.ignoreFile)
		if err != nil {
			return fmt.Errorf("reading ignore file: %w", err)
		}
		opts.Keys = append(opts.Keys, fileOpts.Keys...)
		opts.Prefixes = append(opts.Prefixes, fileOpts.Prefixes...)
	}

	result = diff.ApplyIgnore(result, opts)
	diff.PrintReport(out, result)
	return nil
}

// RunIgnoreMain is the entry point wired from main dispatch.
func RunIgnoreMain(args []string) {
	if err := RunIgnore(args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
