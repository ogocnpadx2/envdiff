package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

type maskArgs struct {
	path        string
	revealChars int
	maskChar    string
	patterns    []string
	report      bool
}

func parseMaskArgs(args []string) (maskArgs, error) {
	if len(args) == 0 {
		return maskArgs{}, fmt.Errorf("usage: envdiff mask <file> [--reveal=N] [--char=*] [--patterns=A,B] [--report]")
	}
	a := maskArgs{
		path:        args[0],
		revealChars: 0,
		maskChar:    "*",
		patterns:    diff.DefaultMaskOptions().Patterns,
	}
	for _, flag := range args[1:] {
		switch {
		case strings.HasPrefix(flag, "--reveal="):
			n, err := strconv.Atoi(strings.TrimPrefix(flag, "--reveal="))
			if err != nil || n < 0 {
				return maskArgs{}, fmt.Errorf("invalid --reveal value: %s", flag)
			}
			a.revealChars = n
		case strings.HasPrefix(flag, "--char="):
			ch := strings.TrimPrefix(flag, "--char=")
			if ch == "" {
				return maskArgs{}, fmt.Errorf("--char must not be empty")
			}
			a.maskChar = ch
		case strings.HasPrefix(flag, "--patterns="):
			raw := strings.TrimPrefix(flag, "--patterns=")
			a.patterns = strings.Split(raw, ",")
		case flag == "--report":
			a.report = true
		default:
			return maskArgs{}, fmt.Errorf("unknown flag: %s", flag)
		}
	}
	return a, nil
}

func RunMask(args []string, out io.Writer) error {
	a, err := parseMaskArgs(args)
	if err != nil {
		return err
	}
	env, err := parser.ParseFile(a.path)
	if err != nil {
		return fmt.Errorf("parse %s: %w", a.path, err)
	}
	opts := diff.MaskOptions{
		MaskChar:    a.maskChar,
		RevealChars: a.revealChars,
		Patterns:    a.patterns,
	}
	masked, keys := diff.MaskEnvWithReport(env, opts)
	if a.report {
		if len(keys) == 0 {
			fmt.Fprintln(out, "No keys masked.")
		} else {
			fmt.Fprintf(out, "Masked keys (%d):\n", len(keys))
			for _, k := range keys {
				fmt.Fprintf(out, "  %s\n", k)
			}
		}
		return nil
	}
	allKeys := make([]string, 0, len(masked))
	for k := range masked {
		allKeys = append(allKeys, k)
	}
	sort.Strings(allKeys)
	for _, k := range allKeys {
		fmt.Fprintf(out, "%s=%s\n", k, masked[k])
	}
	return nil
}

func RunMaskMain(args []string) {
	if err := RunMask(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
