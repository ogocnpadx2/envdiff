package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

type renameArgs struct {
	leftPath  string
	rightPath string
	renames   diff.RenameMap
	jsonOut   bool
}

func parseRenameArgs(args []string) (renameArgs, error) {
	fs := flag.NewFlagSet("rename", flag.ContinueOnError)
	rawRenames := fs.String("rename", "", "comma-separated old=new key pairs, e.g. DB_HOST=DATABASE_HOST,FOO=BAR")
	jsonOut := fs.Bool("json", false, "output as JSON")

	if err := fs.Parse(args); err != nil {
		return renameArgs{}, err
	}
	positional := fs.Args()
	if len(positional) < 2 {
		return renameArgs{}, fmt.Errorf("usage: envdiff rename [flags] <left> <right>")
	}

	rm := diff.RenameMap{}
	if *rawRenames != "" {
		for _, pair := range strings.Split(*rawRenames, ",") {
			parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
			if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
				return renameArgs{}, fmt.Errorf("invalid rename pair: %q", pair)
			}
			rm[parts[0]] = parts[1]
		}
	}

	return renameArgs{
		leftPath:  positional[0],
		rightPath: positional[1],
		renames:   rm,
		jsonOut:   *jsonOut,
	}, nil
}

func RunRenameMain(args []string, out io.Writer) error {
	ra, err := parseRenameArgs(args)
	if err != nil {
		return err
	}

	left, err := parser.ParseFile(ra.leftPath)
	if err != nil {
		return fmt.Errorf("parsing left: %w", err)
	}
	right, err := parser.ParseFile(ra.rightPath)
	if err != nil {
		return fmt.Errorf("parsing right: %w", err)
	}

	result := diff.Compare(left, right)
	updated, rr := diff.ApplyRenames(result, ra.renames)

	if ra.jsonOut {
		payload := map[string]interface{}{
			"applied":         rr.Applied,
			"skipped":         rr.Skipped,
			"conflicts":       rr.Conflicts,
			"remaining_diff":  updated,
		}
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(payload)
	}

	fmt.Fprintf(out, "Applied renames: %d\n", len(rr.Applied))
	for _, e := range rr.Applied {
		fmt.Fprintf(out, "  %s -> %s\n", e.OldKey, e.NewKey)
	}
	if len(rr.Skipped) > 0 {
		fmt.Fprintf(out, "Skipped renames: %d\n", len(rr.Skipped))
		for _, e := range rr.Skipped {
			fmt.Fprintf(out, "  %s -> %s (not found in both)\n", e.OldKey, e.NewKey)
		}
	}
	if len(rr.Conflicts) > 0 {
		fmt.Fprintf(out, "Conflicts: %d\n", len(rr.Conflicts))
		for _, e := range rr.Conflicts {
			fmt.Fprintf(out, "  %s -> %s (both keys present)\n", e.OldKey, e.NewKey)
		}
	}
	diff.PrintReport(out, updated)
	return nil
}

func RunRename(args []string) {
	if err := RunRenameMain(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
