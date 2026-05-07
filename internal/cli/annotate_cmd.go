package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

type annotateArgs struct {
	leftPath  string
	rightPath string
	noEmpty   bool
	noURL     bool
	noPlaceholder bool
}

func parseAnnotateArgs(args []string) (annotateArgs, error) {
	fs := flag.NewFlagSet("annotate", flag.ContinueOnError)
	noEmpty := fs.Bool("no-empty", false, "disable empty-value annotations")
	noURL := fs.Bool("no-url", false, "disable URL-value annotations")
	noPlaceholder := fs.Bool("no-placeholder", false, "disable placeholder-value annotations")

	if err := fs.Parse(args); err != nil {
		return annotateArgs{}, err
	}
	positional := fs.Args()
	if len(positional) < 2 {
		return annotateArgs{}, fmt.Errorf("usage: envdiff annotate [options] <left.env> <right.env>")
	}
	return annotateArgs{
		leftPath:      positional[0],
		rightPath:     positional[1],
		noEmpty:       *noEmpty,
		noURL:         *noURL,
		noPlaceholder: *noPlaceholder,
	}, nil
}

// RunAnnotate executes the annotate command writing output to w.
func RunAnnotate(args []string, w io.Writer) error {
	a, err := parseAnnotateArgs(args)
	if err != nil {
		return err
	}

	leftEnv, err := parser.ParseFile(a.leftPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", a.leftPath, err)
	}
	rightEnv, err := parser.ParseFile(a.rightPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", a.rightPath, err)
	}

	r := diff.Compare(leftEnv, rightEnv)
	opts := diff.AnnotateOptions{
		NoteEmpty:       !a.noEmpty,
		NoteURL:         !a.noURL,
		NotePlaceholder: !a.noPlaceholder,
	}
	ar := diff.Annotate(r, leftEnv, rightEnv, opts)

	if len(ar.Annotations) == 0 {
		fmt.Fprintln(w, "No annotations.")
		return nil
	}
	for _, note := range ar.Annotations {
		fmt.Fprintf(w, "  [%s] %s\n", note.Key, note.Message)
	}
	return nil
}

// RunAnnotateMain is the entry-point called from main when the subcommand is "annotate".
func RunAnnotateMain(args []string) {
	if err := RunAnnotate(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
