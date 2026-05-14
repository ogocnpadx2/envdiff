package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

type patchArgs struct {
	srcPath string
	dstPath string
	targetPath string
	dryRun bool
	show bool
}

func parsePatchArgs(args []string) (patchArgs, error) {
	fs := flag.NewFlagSet("patch", flag.ContinueOnError)
	dryRun := fs.Bool("dry-run", false, "preview changes without applying")
	show := fs.Bool("show", false, "print the patch in unified format and exit")

	if err := fs.Parse(args); err != nil {
		return patchArgs{}, err
	}

	remaining := fs.Args()
	if len(remaining) < 2 {
		return patchArgs{}, fmt.Errorf("usage: envdiff patch [--dry-run] [--show] <src> <dst> [target]")
	}

	pa := patchArgs{
		srcPath: remaining[0],
		dstPath: remaining[1],
		dryRun:  *dryRun,
		show:    *show,
	}
	if len(remaining) >= 3 {
		pa.targetPath = remaining[2]
	}
	return pa, nil
}

// RunPatch executes the patch subcommand writing output to w.
func RunPatch(args []string, w io.Writer) error {
	pa, err := parsePatchArgs(args)
	if err != nil {
		return err
	}

	src, err := parser.ParseFile(pa.srcPath)
	if err != nil {
		return fmt.Errorf("reading src: %w", err)
	}
	dst, err := parser.ParseFile(pa.dstPath)
	if err != nil {
		return fmt.Errorf("reading dst: %w", err)
	}

	entries := diff.BuildPatch(src, dst)

	if pa.show || pa.targetPath == "" {
		fmt.Fprint(w, diff.FormatPatch(entries))
		return nil
	}

	target, err := parser.ParseFile(pa.targetPath)
	if err != nil {
		return fmt.Errorf("reading target: %w", err)
	}

	_, result := diff.ApplyPatch(target, entries, pa.dryRun)

	fmt.Fprintf(w, "Applied: %d  Skipped: %d  Conflicts: %d\n",
		len(result.Applied), len(result.Skipped), len(result.Conflicts))

	if len(result.Conflicts) > 0 {
		fmt.Fprintln(w, "Conflicts:")
		for _, c := range result.Conflicts {
			fmt.Fprintf(w, "  %s (expected %q, found different value)\n", c.Key, c.OldValue)
		}
	}
	return nil
}

// RunPatchMain is the entry point called from main.
func RunPatchMain(args []string) {
	if err := RunPatch(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
