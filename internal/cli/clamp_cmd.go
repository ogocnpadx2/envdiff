package cli

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

type clampArgs struct {
	path   string
	min    *float64
	max    *float64
	keys   []string
	strict bool
}

func parseClampArgs(args []string) (clampArgs, error) {
	if len(args) < 1 {
		return clampArgs{}, fmt.Errorf("usage: envdiff clamp <file> [--min N] [--max N] [--keys K1,K2] [--strict]")
	}
	a := clampArgs{path: args[0]}
	for i := 1; i < len(args); i++ {
		switch {
		case args[i] == "--strict":
			a.strict = true
		case strings.HasPrefix(args[i], "--min="):
			v, err := strconv.ParseFloat(strings.TrimPrefix(args[i], "--min="), 64)
			if err != nil {
				return clampArgs{}, fmt.Errorf("invalid --min value: %w", err)
			}
			a.min = &v
		case strings.HasPrefix(args[i], "--max="):
			v, err := strconv.ParseFloat(strings.TrimPrefix(args[i], "--max="), 64)
			if err != nil {
				return clampArgs{}, fmt.Errorf("invalid --max value: %w", err)
			}
			a.max = &v
		case strings.HasPrefix(args[i], "--keys="):
			raw := strings.TrimPrefix(args[i], "--keys=")
			for _, k := range strings.Split(raw, ",") {
				if k = strings.TrimSpace(k); k != "" {
					a.keys = append(a.keys, k)
				}
			}
		default:
			return clampArgs{}, fmt.Errorf("unknown flag: %s", args[i])
		}
	}
	return a, nil
}

func RunClamp(args []string, out io.Writer) error {
	a, err := parseClampArgs(args)
	if err != nil {
		return err
	}
	env, err := parser.ParseFile(a.path)
	if err != nil {
		return fmt.Errorf("parse %s: %w", a.path, err)
	}
	opts := diff.ClampOptions{
		Min:    a.min,
		Max:    a.max,
		Keys:   a.keys,
		Strict: a.strict,
	}
	report, err := diff.ClampEnv(env, opts)
	if err != nil {
		return err
	}
	for _, v := range report.Violations {
		fmt.Fprintf(out, "CLAMPED  %s: %s -> %s (%s)\n", v.Key, v.Original, v.Result, v.Reason)
	}
	if len(report.Violations) == 0 {
		fmt.Fprintln(out, "all values within bounds")
	}
	return nil
}

func RunClampMain(args []string) {
	if err := RunClamp(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
