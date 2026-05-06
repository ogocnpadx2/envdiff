package diff

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// ExportFormat represents supported export formats beyond the standard output.
type ExportFormat string

const (
	ExportCSV      ExportFormat = "csv"
	ExportDotEnv   ExportFormat = "dotenv"
)

// ParseExportFormat parses a string into an ExportFormat.
func ParseExportFormat(s string) (ExportFormat, error) {
	switch strings.ToLower(s) {
	case "csv":
		return ExportCSV, nil
	case "dotenv":
		return ExportDotEnv, nil
	default:
		return "", fmt.Errorf("unknown export format %q: must be one of csv, dotenv", s)
	}
}

// Export writes the diff result to w in the given ExportFormat.
func Export(w io.Writer, result Result, format ExportFormat) error {
	switch format {
	case ExportCSV:
		return exportCSV(w, result)
	case ExportDotEnv:
		return exportDotEnv(w, result)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

func exportCSV(w io.Writer, result Result) error {
	csvW := csv.NewWriter(w)
	if err := csvW.Write([]string{"key", "status", "left_value", "right_value"}); err != nil {
		return err
	}
	for _, k := range result.MissingInRight {
		if err := csvW.Write([]string{k, "missing_in_right", "", ""}); err != nil {
			return err
		}
	}
	for _, k := range result.MissingInLeft {
		if err := csvW.Write([]string{k, "missing_in_left", "", ""}); err != nil {
			return err
		}
	}
	for _, m := range result.Mismatched {
		if err := csvW.Write([]string{m.Key, "mismatched", m.LeftVal, m.RightVal}); err != nil {
			return err
		}
	}
	csvW.Flush()
	return csvW.Error()
}

func exportDotEnv(w io.Writer, result Result) error {
	if len(result.MissingInRight) > 0 {
		fmt.Fprintln(w, "# Keys missing in right file")
		for _, k := range result.MissingInRight {
			fmt.Fprintf(w, "# %s=\n", k)
		}
	}
	if len(result.MissingInLeft) > 0 {
		fmt.Fprintln(w, "# Keys missing in left file")
		for _, k := range result.MissingInLeft {
			fmt.Fprintf(w, "# %s=\n", k)
		}
	}
	if len(result.Mismatched) > 0 {
		fmt.Fprintln(w, "# Mismatched keys (left values shown)")
		for _, m := range result.Mismatched {
			fmt.Fprintf(w, "%s=%s\n", m.Key, m.LeftVal)
		}
	}
	return nil
}
