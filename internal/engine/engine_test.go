package engine

import (
	"testing"

	"github.com/lucasdeprit/agent-policy-kit/internal/diff"
	"github.com/lucasdeprit/agent-policy-kit/internal/rules"
)

func TestReviewAppliesRegexToIncludedAddedLines(t *testing.T) {
	changes := []diff.Change{
		{File: "src/app.ts", Line: 42, Text: `console.log("debug")`},
		{File: "src/generated/client.ts", Line: 10, Text: `console.log("Allowed")`},
		{File: "test/app.test.ts", Line: 5, Text: `console.log("Ignored")`},
	}
	ruleSet := []rules.Rule{
		{
			ID:       "GEN001",
			Severity: rules.SeverityError,
			Include:  []string{"src/**/*.ts"},
			Exclude:  []string{"src/generated/**"},
			Match:    rules.Match{Type: "regex", Pattern: `\bconsole\.log\s*\(`},
		},
	}

	violations := Review(changes, ruleSet)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].RuleID != "GEN001" || violations[0].File != "src/app.ts" || violations[0].Line != 42 {
		t.Fatalf("unexpected violation: %#v", violations[0])
	}
}
