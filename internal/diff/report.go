package diff

import (
	"fmt"
	"io"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
	colorBold   = "\033[1m"
)

// PrintReport writes a human-readable diff report to w.
func PrintReport(w io.Writer, leftName, rightName string, r Report) {
	if r.Clean() {
		fprintf(w, "%s✔ No differences found between %s and %s%s\n",
			colorGreen, leftName, rightName, colorReset)
		return
	}

	fprintf(w, "%s%senvdiff: %s ↔ %s%s\n\n", colorBold, colorReset, leftName, rightName, colorReset)

	if len(r.MissingInRight) > 0 {
		fprintf(w, "%s%s Missing in %s:%s\n", colorRed, colorBold, rightName, colorReset)
		for _, k := range r.MissingInRight {
			fprintf(w, "  %s- %s%s\n", colorRed, k, colorReset)
		}
		fprintf(w, "\n")
	}

	if len(r.MissingInLeft) > 0 {
		fprintf(w, "%s%s Missing in %s:%s\n", colorYellow, colorBold, leftName, colorReset)
		for _, k := range r.MissingInLeft {
			fprintf(w, "  %s+ %s%s\n", colorYellow, k, colorReset)
		}
		fprintf(w, "\n")
	}

	if len(r.Mismatched) > 0 {
		fprintf(w, "%s%s Value mismatches:%s\n", colorBold, colorYellow, colorReset)
		for _, m := range r.Mismatched {
			fprintf(w, "  %s~ %s%s\n", colorYellow, m.Key, colorReset)
			fprintf(w, "      %s: %q\n", leftName, m.LeftVal)
			fprintf(w, "      %s: %q\n", rightName, m.RightVal)
		}
		fprintf(w, "\n")
	}
}

func fprintf(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, format, args...)
}
