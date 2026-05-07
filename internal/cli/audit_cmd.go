package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

type auditArgs struct {
	leftPath  string
	rightPath string
	jsonOut   bool
}

func parseAuditArgs(args []string) (auditArgs, error) {
	var a auditArgs
	positional := []string{}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--json":
			a.jsonOut = true
		default:
			positional = append(positional, args[i])
		}
	}

	if len(positional) < 2 {
		return a, fmt.Errorf("usage: envdiff audit <left.env> <right.env> [--json]")
	}
	a.leftPath = positional[0]
	a.rightPath = positional[1]
	return a, nil
}

// RunAudit runs the audit subcommand writing output to w.
func RunAudit(args []string, w io.Writer) error {
	a, err := parseAuditArgs(args)
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

	log := diff.BuildAuditLog(a.leftPath, a.rightPath, left, right)

	if a.jsonOut {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(log)
	}

	fmt.Fprintf(w, "Audit: %s vs %s\n", a.leftPath, a.rightPath)
	fmt.Fprintf(w, "Summary: %s\n\n", log.Summary())
	for _, e := range log.Events {
		switch e.Type {
		case diff.AuditAdded:
			fmt.Fprintf(w, "  + %-30s = %s\n", e.Key, e.NewValue)
		case diff.AuditRemoved:
			fmt.Fprintf(w, "  - %-30s (was: %s)\n", e.Key, e.OldValue)
		case diff.AuditChanged:
			fmt.Fprintf(w, "  ~ %-30s %s -> %s\n", e.Key, e.OldValue, e.NewValue)
		case diff.AuditUnchanged:
			fmt.Fprintf(w, "    %-30s (unchanged)\n", e.Key)
		}
	}
	return nil
}

// RunAuditMain is the entry point called from main dispatch.
func RunAuditMain(args []string) {
	if err := RunAudit(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
