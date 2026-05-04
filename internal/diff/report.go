package diff

import (
	"fmt"
	"io"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
)

// PrintReport writes a human-readable diff report to w.
// leftName and rightName are labels for the two env files being compared.
func PrintReport(w io.Writer, result Result, leftName, rightName string, color bool) {
	if result.IsClean() {
		fprintf(w, color, colorGreen, "✔ No differences found between %s and %s\n", leftName, rightName)
		return
	}

	if len(result.MissingInRight) > 0 {
		fprintf(w, color, colorRed, "Keys in %s but missing in %s:\n", leftName, rightName)
		for _, k := range result.MissingInRight {
			fprintf(w, color, colorRed, "  - %s\n", k)
		}
	}

	if len(result.MissingInLeft) > 0 {
		fprintf(w, color, colorRed, "Keys in %s but missing in %s:\n", rightName, leftName)
		for _, k := range result.MissingInLeft {
			fprintf(w, color, colorRed, "  - %s\n", k)
		}
	}

	if len(result.Mismatched) > 0 {
		fprintf(w, color, colorYellow, "Mismatched values:\n")
		for _, m := range result.Mismatched {
			line := fmt.Sprintf("  ~ %s\n      %s: %s\n      %s: %s\n",
				m.Key,
				leftName, m.LeftValue,
				rightName, m.RightValue,
			)
			fprintf(w, color, colorYellow, "%s", line)
		}
	}
}

func fprintf(w io.Writer, useColor bool, clr, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if useColor && !strings.HasPrefix(clr, "") {
		msg = clr + msg + colorReset
	} else if useColor {
		msg = clr + msg + colorReset
	}
	fmt.Fprint(w, msg)
}
