package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

type templateArgs struct {
	leftFile     string
	rightFile    string
	templatePath string
	templateStr  string
}

func parseTemplateArgs(args []string) (templateArgs, error) {
	var a templateArgs
	remaining := []string{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--template-file":
			if i+1 >= len(args) {
				return a, fmt.Errorf("--template-file requires a value")
			}
			i++
			a.templatePath = args[i]
		case "--template":
			if i+1 >= len(args) {
				return a, fmt.Errorf("--template requires a value")
			}
			i++
			a.templateStr = args[i]
		default:
			remaining = append(remaining, args[i])
		}
	}
	if len(remaining) < 2 {
		return a, fmt.Errorf("template command requires two .env file paths")
	}
	if a.templatePath == "" && a.templateStr == "" {
		return a, fmt.Errorf("either --template or --template-file is required")
	}
	a.leftFile = remaining[0]
	a.rightFile = remaining[1]
	return a, nil
}

// RunTemplate parses two env files and renders the user-supplied template.
func RunTemplate(args []string, w io.Writer) error {
	a, err := parseTemplateArgs(args)
	if err != nil {
		return err
	}
	left, err := parser.ParseFile(a.leftFile)
	if err != nil {
		return fmt.Errorf("parse %s: %w", a.leftFile, err)
	}
	right, err := parser.ParseFile(a.rightFile)
	if err != nil {
		return fmt.Errorf("parse %s: %w", a.rightFile, err)
	}
	result := diff.Compare(left, right)
	if a.templatePath != "" {
		return diff.RenderTemplateFile(w, a.templatePath, result, a.leftFile, a.rightFile)
	}
	return diff.RenderTemplate(w, a.templateStr, result, a.leftFile, a.rightFile)
}

// RunTemplateMain is the entry point called from main for the template subcommand.
func RunTemplateMain(args []string) {
	if err := RunTemplate(args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
