package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseOptionsDefaults(t *testing.T) {
	options, err := parseOptions("review", nil)
	if err != nil {
		t.Fatalf("parseOptions returned error: %v", err)
	}
	if options.Target == "" || len(options.RuleSources) != 1 {
		t.Fatalf("unexpected options: %#v", options)
	}
	if filepath.Base(options.RuleSources[0]) != "rules" {
		t.Fatalf("expected default rules source, got %#v", options.RuleSources)
	}
}

func TestParseOptionsTargetAndRepeatedRules(t *testing.T) {
	options, err := parseOptions("review", []string{"--target", "../client", "--rules", "../kit/rules", "--rules", "./rules"})
	if err != nil {
		t.Fatalf("parseOptions returned error: %v", err)
	}
	if filepath.Base(options.Target) != "client" {
		t.Fatalf("expected target client, got %s", options.Target)
	}
	if len(options.RuleSources) != 2 {
		t.Fatalf("expected 2 rules sources, got %#v", options.RuleSources)
	}
}

func TestInitProjectCreatesRules(t *testing.T) {
	dir := t.TempDir()

	created, skipped, err := initProject(dir)
	if err != nil {
		t.Fatalf("initProject returned error: %v", err)
	}
	if created != 2 || skipped != 0 {
		t.Fatalf("expected created=2 skipped=0, got created=%d skipped=%d", created, skipped)
	}

	expected := []string{
		filepath.Join(dir, "rules", "default.yaml"),
		filepath.Join(dir, "rules", "examples", "rule-template.yaml"),
	}
	for _, path := range expected {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}
}

func TestInitProjectDoesNotOverwriteExistingRules(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rules", "default.yaml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte("custom: true\n"), 0o644); err != nil {
		t.Fatalf("write custom rule: %v", err)
	}

	created, skipped, err := initProject(dir)
	if err != nil {
		t.Fatalf("initProject returned error: %v", err)
	}
	if created != 1 || skipped != 1 {
		t.Fatalf("expected created=1 skipped=1, got created=%d skipped=%d", created, skipped)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read custom rule: %v", err)
	}
	if string(content) != "custom: true\n" {
		t.Fatalf("existing rule was overwritten: %q", string(content))
	}
}
