package parser

import (
	"os"
	"path/filepath"
	"testing"
)

const testRule = `rules:
  - id: T001
    name: no_foo
    severity: error
    include:
      - "**/*.txt"
    match:
      type: regex
      pattern: "foo"
    message: "Do not use foo."
`

func TestLoadRuleSourcesLoadsMultipleSources(t *testing.T) {
	dir := t.TempDir()
	rulesA := filepath.Join(dir, "rules-a")
	rulesB := filepath.Join(dir, "rules-b")
	if err := os.MkdirAll(rulesA, 0o755); err != nil {
		t.Fatalf("mkdir rules-a: %v", err)
	}
	if err := os.MkdirAll(rulesB, 0o755); err != nil {
		t.Fatalf("mkdir rules-b: %v", err)
	}
	if err := os.WriteFile(filepath.Join(rulesA, "a.yaml"), []byte(testRule), 0o644); err != nil {
		t.Fatalf("write a rule: %v", err)
	}
	if err := os.WriteFile(filepath.Join(rulesB, "b.yml"), []byte(testRule), 0o644); err != nil {
		t.Fatalf("write b rule: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(rulesA, "examples"), 0o755); err != nil {
		t.Fatalf("mkdir examples: %v", err)
	}
	if err := os.WriteFile(filepath.Join(rulesA, "examples", "ignored.yaml"), []byte(testRule), 0o644); err != nil {
		t.Fatalf("write ignored rule: %v", err)
	}

	rules, err := LoadRuleSources([]string{rulesA, rulesB})
	if err != nil {
		t.Fatalf("LoadRuleSources returned error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 loaded rules, got %d", len(rules))
	}
}

func TestRuleFilesFromSourceSupportsGlob(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rule.yaml")
	if err := os.WriteFile(path, []byte(testRule), 0o644); err != nil {
		t.Fatalf("write rule: %v", err)
	}

	files, err := RuleFilesFromSource(filepath.Join(dir, "*.yaml"))
	if err != nil {
		t.Fatalf("RuleFilesFromSource returned error: %v", err)
	}
	if len(files) != 1 || files[0] != path {
		t.Fatalf("unexpected files: %#v", files)
	}
}
