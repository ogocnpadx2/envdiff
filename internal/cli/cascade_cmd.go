package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

type cascadeArgs struct {
	paths    []string
	labels   []string
	strategy string
	showSkip bool
}

func parseCascadeArgs(args []string) (cascadeArgs, error) {
	var ca cascadeArgs
	ca.strategy = "overwrite"
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--strategy":
			if i+1 >= len(args) {
				return ca, fmt.Errorf("--strategy requires a value")
			}
			i++
			ca.strategy = args[i]
		case "--labels":
			if i+1 >= len(args) {
				return ca, fmt.Errorf("--labels requires a value")
			}
			i++
			ca.labels = strings.Split(args[i], ",")
		case "--show-skipped":
			ca.showSkip = true
		default:
			ca.paths = append(ca.paths, args[i])
		}
	}
	if len(ca.paths) < 2 {
		return ca, fmt.Errorf("cascade requires at least two .env file paths")
	}
	return ca, nil
}

// RunCascade executes the cascade sub-command.
func RunCascade(args []string, out io.Writer) error {
	ca, err := parseCascadeArgs(args)
	if err != nil {
		return err
	}
	strategy, err := diff.ParseCascadeStrategy(ca.strategy)
	if err != nil {
		return err
	}
	envs := make([]map[string]string, 0, len(ca.paths))
	labels := ca.labels
	for i, p := range ca.paths {
		env, err := parser.ParseFile(p)
		if err != nil {
			return fmt.Errorf("reading %s: %w", p, err)
		}
		envs = append(envs, env)
		if i >= len(labels) {
			labels = append(labels, p)
		}
	}
	res := diff.Cascade(envs, labels, strategy)
	fmt.Fprintf(out, "# cascaded result (%s strategy)\n", ca.strategy)
	for _, e := range res.Resolved {
		fmt.Fprintf(out, "%s=%s  # from %s\n", e.Key, e.Value, e.Origin)
	}
	if ca.showSkip && len(res.Skipped) > 0 {
		fmt.Fprintln(out, "\n# skipped (shadowed) entries")
		for _, e := range res.Skipped {
			fmt.Fprintf(out, "# %s=%s  (from %s)\n", e.Key, e.Value, e.Origin)
		}
	}
	return nil
}

// RunCascadeMain is the entry point called from main dispatch.
func RunCascadeMain(args []string) {
	if err := RunCascade(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
