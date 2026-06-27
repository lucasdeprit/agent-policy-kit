package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestDoctorReportsOKForGitTargetAndRules(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	rulesDir := filepath.Join(dir, "rules")
	mkdir(t, target)
	mkdir(t, rulesDir)
	runGit(t, target, "init")
	write(t, filepath.Join(rulesDir, "rules.yaml"), `rules:
  - id: T001
    name: no_foo
    severity: error
    include:
      - "**/*.txt"
    match:
      type: regex
      pattern: "foo"
    message: "Do not use foo."
`)

	report := Doctor(Options{Target: target, RuleSources: []string{rulesDir}})
	if len(report.Failures) != 0 {
		t.Fatalf("expected no failures, got %#v", report.Failures)
	}
	if len(report.RuleFiles) != 1 {
		t.Fatalf("expected one rule file, got %#v", report.RuleFiles)
	}
}

func TestDoctorFailsForMissingRules(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	mkdir(t, target)
	runGit(t, target, "init")

	report := Doctor(Options{Target: target, RuleSources: []string{filepath.Join(dir, "missing")}})
	if len(report.Failures) == 0 {
		t.Fatalf("expected failures for missing rules")
	}
}

func mkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, string(out))
	}
}
