package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

type encryptArgs struct {
	path     string
	hashKeys []string
	prefix   string
	showReport bool
}

func parseEncryptArgs(args []string) (encryptArgs, error) {
	if len(args) == 0 {
		return encryptArgs{}, fmt.Errorf("usage: envdiff encrypt <file> [--hash-keys=k1,k2] [--prefix=sha256:] [--report]")
	}
	a := encryptArgs{
		path:   args[0],
		prefix: "sha256:",
	}
	opts := diff.DefaultEncryptOptions()
	a.hashKeys = opts.HashKeys

	for _, flag := range args[1:] {
		switch {
		case strings.HasPrefix(flag, "--hash-keys="):
			val := strings.TrimPrefix(flag, "--hash-keys=")
			a.hashKeys = strings.Split(val, ",")
		case strings.HasPrefix(flag, "--prefix="):
			a.prefix = strings.TrimPrefix(flag, "--prefix=")
		case flag == "--report":
			a.showReport = true
		default:
			return encryptArgs{}, fmt.Errorf("unknown flag: %s", flag)
		}
	}
	return a, nil
}

// RunEncrypt parses an env file, hashes sensitive values, and prints the result.
func RunEncrypt(args []string, out io.Writer) error {
	a, err := parseEncryptArgs(args)
	if err != nil {
		return err
	}

	env, err := parser.ParseFile(a.path)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	opts := diff.EncryptOptions{
		HashKeys: a.hashKeys,
		Prefix:   a.prefix,
	}

	encrypted, report := diff.EncryptEnvWithReport(env, opts)

	if a.showReport {
		if len(report.HashedKeys) == 0 {
			fmt.Fprintln(out, "# no keys hashed")
		} else {
			fmt.Fprintf(out, "# hashed keys: %s\n", strings.Join(report.HashedKeys, ", "))
		}
	}

	for _, k := range sortedKeys(encrypted) {
		fmt.Fprintf(out, "%s=%s\n", k, encrypted[k])
	}
	return nil
}

// RunEncryptMain is the entry point wired into main dispatch.
func RunEncryptMain(args []string) {
	if err := RunEncrypt(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// sortedKeys returns map keys in sorted order (reuses existing helper if available).
func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	diff.SortStringsExported(keys)
	return keys
}
