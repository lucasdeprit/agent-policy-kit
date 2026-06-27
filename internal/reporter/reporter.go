package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/lucasdeprit/agent-policy-kit/internal/engine"
)

func Text(w io.Writer, violations []engine.Violation) {
	if len(violations) == 0 {
		fmt.Fprintln(w, "OK")
		return
	}

	fmt.Fprintf(w, "FAILED %d\n", len(violations))
	for _, violation := range violations {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "%s %s:%d\n", violation.RuleID, violation.File, violation.Line)
		fmt.Fprintf(w, "+ %s\n", strings.TrimSpace(violation.Text))
		if violation.Expected != "" {
			fmt.Fprintf(w, "expected: %s\n", violation.Expected)
		}
		if violation.Suggestion != "" {
			fmt.Fprintf(w, "fix: %s\n", violation.Suggestion)
		}
	}
}
