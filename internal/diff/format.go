package diff

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// OutputFormat represents the output format for the diff report.
type OutputFormat int

const (
	FormatText OutputFormat = iota
	FormatJSON
	FormatMarkdown
)

// ParseFormat parses a format string into an OutputFormat.
func ParseFormat(s string) (OutputFormat, error) {
	switch strings.ToLower(s) {
	case "", "text":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	case "markdown", "md":
		return FormatMarkdown, nil
	default:
		return FormatText, fmt.Errorf("unknown format %q: must be text, json, or markdown", s)
	}
}

// jsonReport is the structure used for JSON output.
type jsonReport struct {
	Clean      bool        `json:"clean"`
	MissingIn  missingJSON `json:"missing"`
	Mismatched []mismatchJSON `json:"mismatched"`
}

type missingJSON struct {
	Left  []string `json:"left"`
	Right []string `json:"right"`
}

type mismatchJSON struct {
	Key        string `json:"key"`
	LeftValue  string `json:"left_value"`
	RightValue string `json:"right_value"`
}

// PrintFormatted writes the diff result in the requested format.
func PrintFormatted(w io.Writer, result Result, leftName, rightName string, format OutputFormat) error {
	switch format {
	case FormatJSON:
		return printJSON(w, result)
	case FormatMarkdown:
		return printMarkdown(w, result, leftName, rightName)
	default:
		return printText(w, result, leftName, rightName)
	}
}

func printJSON(w io.Writer, result Result) error {
	var mismatches []mismatchJSON
	for _, m := range result.Mismatched {
		mismatches = append(mismatches, mismatchJSON{
			Key:        m.Key,
			LeftValue:  m.LeftValue,
			RightValue: m.RightValue,
		})
	}
	report := jsonReport{
		Clean: result.IsClean(),
		MissingIn: missingJSON{
			Left:  nullSafe(result.MissingInLeft),
			Right: nullSafe(result.MissingInRight),
		},
		Mismatched: mismatches,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

func printMarkdown(w io.Writer, result Result, leftName, rightName string) error {
	fmt.Fprintf(w, "## envdiff: `%s` vs `%s`\n\n", leftName, rightName)
	if result.IsClean() {
		fmt.Fprintln(w, "✅ No differences found.")
		return nil
	}
	if len(result.MissingInRight) > 0 {
		fmt.Fprintf(w, "### Missing in `%s`\n\n", rightName)
		for _, k := range result.MissingInRight {
			fmt.Fprintf(w, "- `%s`\n", k)
		}
		fmt.Fprintln(w)
	}
	if len(result.MissingInLeft) > 0 {
		fmt.Fprintf(w, "### Missing in `%s`\n\n", leftName)
		for _, k := range result.MissingInLeft {
			fmt.Fprintf(w, "- `%s`\n", k)
		}
		fmt.Fprintln(w)
	}
	if len(result.Mismatched) > 0 {
		fmt.Fprintln(w, "### Mismatched values\n")
		fmt.Fprintf(w, "| Key | `%s` | `%s` |\n", leftName, rightName)
		fmt.Fprintln(w, "|-----|------|-------|")
		for _, m := range result.Mismatched {
			fmt.Fprintf(w, "| `%s` | `%s` | `%s` |\n", m.Key, m.LeftValue, m.RightValue)
		}
	}
	return nil
}

func printText(w io.Writer, result Result, leftName, rightName string) error {
	if result.IsClean() {
		fmt.Fprintln(w, "No differences found.")
		return nil
	}
	for _, k := range result.MissingInRight {
		fmt.Fprintf(w, "MISSING_IN_RIGHT\t%s\n", k)
	}
	for _, k := range result.MissingInLeft {
		fmt.Fprintf(w, "MISSING_IN_LEFT\t%s\n", k)
	}
	for _, m := range result.Mismatched {
		fmt.Fprintf(w, "MISMATCH\t%s\t%s=%s\t%s=%s\n", m.Key, leftName, m.LeftValue, rightName, m.RightValue)
	}
	return nil
}

func nullSafe(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
