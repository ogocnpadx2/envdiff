package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

type aliasArgs struct {
	envPath string
	aliases []string
}

func parseAliasArgs(args []string) (aliasArgs, error) {
	var a aliasArgs
	var remaining []string
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--alias" && i+1 < len(args):
			i++
			a.aliases = append(a.aliases, args[i])
		case strings.HasPrefix(args[i], "--alias="):
			a.aliases = append(a.aliases, strings.TrimPrefix(args[i], "--alias="))
		default:
			remaining = append(remaining, args[i])
		}
	}
	if len(remaining) < 1 {
		return a, fmt.Errorf("usage: envdiff alias <env-file> --alias canonical=alt1,alt2")
	}
	if len(a.aliases) == 0 {
		return a, fmt.Errorf("at least one --alias mapping is required")
	}
	a.envPath = remaining[0]
	return a, nil
}

// RunAlias resolves alias mappings for a single env file and prints a report.
func RunAlias(args []string, out io.Writer) error {
	a, err := parseAliasArgs(args)
	if err != nil {
		return err
	}
	env, err := parser.ParseFile(a.envPath)
	if err != nil {
		return fmt.Errorf("parse %s: %w", a.envPath, err)
	}
	am, err := diff.ParseAliasMap(a.aliases)
	if err != nil {
		return fmt.Errorf("parse aliases: %w", err)
	}
	res := diff.ResolveAliases(env, am)

	if len(res.Resolved) > 0 {
		fmt.Fprintln(out, "Resolved keys:")
		for k, v := range res.Resolved {
			if alias, ok := res.UsedAlias[k]; ok {
				fmt.Fprintf(out, "  %s = %s  (via alias %s)\n", k, v, alias)
			} else {
				fmt.Fprintf(out, "  %s = %s\n", k, v)
			}
		}
	}
	if len(res.Unresolved) > 0 {
		fmt.Fprintln(out, "Unresolved canonical keys:")
		for _, k := range res.Unresolved {
			fmt.Fprintf(out, "  %s\n", k)
		}
	}
	return nil
}

// RunAliasMain is the entry point wired from main dispatch.
func RunAliasMain(args []string) {
	if err := RunAlias(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
