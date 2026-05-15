package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

type truncateArgs struct {
	path      string
	maxLength int
	suffix    string
	keyFilter string
	report    bool
}

func parseTruncateArgs(args []string) (truncateArgs, error) {
	if len(args) < 1 {
		return truncateArgs{}, fmt.Errorf("usage: envdiff truncate <file> [--max=N] [--suffix=S] [--key-filter=F] [--report]")
	}
	a := truncateArgs{
		path:      args[0],
		maxLength: 40,
		suffix:    "...",
	}
	for _, arg := range args[1:] {
		switch {
		case strings.HasPrefix(arg, "--max="):
			n, err := strconv.Atoi(strings.TrimPrefix(arg, "--max="))
			if err != nil || n <= 0 {
				return truncateArgs{}, fmt.Errorf("invalid --max value: %s", arg)
			}
			a.maxLength = n
		case strings.HasPrefix(arg, "--suffix="):
			a.suffix = strings.TrimPrefix(arg, "--suffix=")
		case strings.HasPrefix(arg, "--key-filter="):
			a.keyFilter = strings.TrimPrefix(arg, "--key-filter=")
		case arg == "--report":
			a.report = true
		default:
			return truncateArgs{}, fmt.Errorf("unknown flag: %s", arg)
		}
	}
	return a, nil
}

func RunTruncate(args []string, out io.Writer) error {
	a, err := parseTruncateArgs(args)
	if err != nil {
		return err
	}
	env, err := parser.ParseFile(a.path)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}
	opts := diff.TruncateOptions{
		MaxLength: a.maxLength,
		Suffix:    a.suffix,
		KeyFilter: a.keyFilter,
	}
	result, truncated := diff.TruncateEnvWithReport(env, opts)

	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(out, "%s=%s\n", k, result[k])
	}
	if a.report && len(truncated) > 0 {
		fmt.Fprintf(out, "\n# truncated keys (%d): %s\n", len(truncated), strings.Join(truncated, ", "))
	}
	return nil
}

func RunTruncateMain(args []string) {
	if err := RunTruncate(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
