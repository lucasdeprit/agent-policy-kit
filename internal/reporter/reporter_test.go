package reporter

import (
	"bytes"
	"testing"

	"github.com/lucasdeprit/agent-policy-kit/internal/engine"
)

func TestTextReporter(t *testing.T) {
	var buf bytes.Buffer
	Text(&buf, []engine.Violation{
		{
			RuleID:     "GEN001",
			File:       "src/app.ts",
			Line:       42,
			Text:       ` console.log("debug")`,
			Expected:   "Use the approved abstraction.",
			Suggestion: "Replace the raw API with the project-approved alternative.",
		},
	})

	want := "FAILED 1\n\nGEN001 src/app.ts:42\n+ console.log(\"debug\")\nexpected: Use the approved abstraction.\nfix: Replace the raw API with the project-approved alternative.\n"
	if buf.String() != want {
		t.Fatalf("unexpected output:\n%s", buf.String())
	}
}
