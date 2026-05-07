package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

type transformArgs struct {
	filePath   string
	transforms []string
	output     string // "text" or "dotenv"
}

func parseTransformArgs(args []string) (transformArgs, error) {
	if len(args) < 1 {
		return transformArgs{}, fmt.Errorf("usage: envdiff transform <file> [--transform=<opt>,...] [--output=text|dotenv]")
	}
	a := transformArgs{
		filePath: args[0],
		output:   "text",
	}
	for _, arg := range args[1:] {
		switch {
		case strings.HasPrefix(arg, "--transform="):
			raw := strings.TrimPrefix(arg, "--transform=")
			a.transforms = strings.Split(raw, ",")
		case strings.HasPrefix(arg, "--output="):
			a.output = strings.TrimPrefix(arg, "--output=")
		default:
			return transformArgs{}, fmt.Errorf("unknown flag: %q", arg)
		}
	}
	return a, nil
}

// RunTransform applies transformations to a .env file and prints the result.
func RunTransform(args []string, out io.Writer) error {
	a, err := parseTransformArgs(args)
	if err != nil {
		return err
	}
	env, err := parser.ParseFile(a.filePath)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}
	opts, err := diff.ParseTransformOptions(a.transforms)
	if err != nil {
		return fmt.Errorf("transform option error: %w", err)
	}
	transformed := diff.TransformEnv(env, opts)
	switch a.output {
	case "dotenv":
		keys := make([]string, 0, len(transformed))
		for k := range transformed {
			keys = append(keys, k)
		}
		for _, k := range keys {
			fmt.Fprintf(out, "%s=%s\n", k, transformed[k])
		}
	default:
		for k, v := range transformed {
			fmt.Fprintf(out, "%-30s = %s\n", k, v)
		}
	}
	return nil
}

// RunTransformMain is the entry point for the transform subcommand.
func RunTransformMain(args []string) {
	if err := RunTransform(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
