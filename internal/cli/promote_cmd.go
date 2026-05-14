package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

type promoteArgs struct {
	srcPath  string
	dstPath  string
	strategy string
	dryRun   bool
}

func parsePromoteArgs(args []string) (promoteArgs, error) {
	if len(args) < 2 {
		return promoteArgs{}, fmt.Errorf("usage: envdiff promote <source> <target> [--strategy=missing|overwrite] [--dry-run]")
	}
	a := promoteArgs{
		srcPath:  args[0],
		dstPath:  args[1],
		strategy: "missing",
	}
	for _, flag := range args[2:] {
		switch {
		case strings.HasPrefix(flag, "--strategy="):
			a.strategy = strings.TrimPrefix(flag, "--strategy=")
		case flag == "--dry-run":
			a.dryRun = true
		default:
			return promoteArgs{}, fmt.Errorf("unknown flag: %s", flag)
		}
	}
	return a, nil
}

// RunPromote executes the promote subcommand, writing to out.
func RunPromote(args []string, out io.Writer) error {
	a, err := parsePromoteArgs(args)
	if err != nil {
		return err
	}

	strategy, err := diff.ParsePromoteStrategy(a.strategy)
	if err != nil {
		return err
	}

	src, err := parser.ParseFile(a.srcPath)
	if err != nil {
		return fmt.Errorf("reading source: %w", err)
	}
	dst, err := parser.ParseFile(a.dstPath)
	if err != nil {
		return fmt.Errorf("reading target: %w", err)
	}

	res := diff.Promote(src, dst, diff.PromoteOptions{
		Strategy: strategy,
		DryRun:   a.dryRun,
	})

	if a.dryRun {
		fmt.Fprintf(out, "[dry-run] would add: %s\n", strings.Join(res.Added, ", "))
		fmt.Fprintf(out, "[dry-run] would overwrite: %s\n", strings.Join(res.Overwritten, ", "))
		fmt.Fprintf(out, "[dry-run] would skip: %s\n", strings.Join(res.Skipped, ", "))
		return nil
	}

	fmt.Fprintf(out, "added: %d  overwritten: %d  skipped: %d\n",
		len(res.Added), len(res.Overwritten), len(res.Skipped))
	return nil
}

// RunPromoteMain is the entry point called from main dispatch.
func RunPromoteMain(args []string) {
	if err := RunPromote(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
