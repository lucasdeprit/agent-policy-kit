package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lucasdeprit/agent-policy-kit/internal/rules"
	"gopkg.in/yaml.v3"
)

func LoadRulesDir(dir string) ([]rules.Rule, error) {
	files, err := RuleFilesFromSource(dir)
	if err != nil {
		return nil, err
	}
	return loadRuleFiles(files, dir)
}

func LoadRuleSources(sources []string) ([]rules.Rule, error) {
	if len(sources) == 0 {
		return nil, fmt.Errorf("at least one rules source is required")
	}

	var files []string
	for _, source := range sources {
		sourceFiles, err := RuleFilesFromSource(source)
		if err != nil {
			return nil, err
		}
		files = append(files, sourceFiles...)
	}

	return loadRuleFiles(files, strings.Join(sources, ", "))
}

func RuleFilesFromSource(source string) ([]string, error) {
	if strings.ContainsAny(source, "*?[") {
		matches, err := filepath.Glob(source)
		if err != nil {
			return nil, fmt.Errorf("invalid rules glob %q: %w", source, err)
		}
		return yamlFiles(matches), nil
	}

	info, err := os.Stat(source)
	if err != nil {
		return nil, fmt.Errorf("read rules source %q: %w", source, err)
	}
	if !info.IsDir() {
		if !isYAML(source) {
			return nil, fmt.Errorf("rules file %q must be .yaml or .yml", source)
		}
		return []string{source}, nil
	}

	entries, err := os.ReadDir(source)
	if err != nil {
		return nil, fmt.Errorf("read rules directory %q: %w", source, err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() || !isYAML(entry.Name()) {
			continue
		}
		files = append(files, filepath.Join(source, entry.Name()))
	}
	return files, nil
}

func loadRuleFiles(files []string, sourceLabel string) ([]rules.Rule, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no rules found in %q", sourceLabel)
	}

	var loaded []rules.Rule
	for _, path := range files {
		fileRules, err := LoadRulesFile(path)
		if err != nil {
			return nil, err
		}
		loaded = append(loaded, fileRules...)
	}

	return loaded, nil
}

func LoadRulesFile(path string) ([]rules.Rule, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read rules file %q: %w", path, err)
	}

	var parsed rules.File
	if err := yaml.Unmarshal(content, &parsed); err != nil {
		return nil, fmt.Errorf("parse rules file %q: %w", path, err)
	}

	for i, rule := range parsed.Rules {
		if err := validateRule(rule); err != nil {
			return nil, fmt.Errorf("%s rule %d: %w", path, i+1, err)
		}
	}

	return parsed.Rules, nil
}

func validateRule(rule rules.Rule) error {
	if rule.ID == "" {
		return fmt.Errorf("id is required")
	}
	if rule.Severity != rules.SeverityError {
		return fmt.Errorf("unsupported severity %q", rule.Severity)
	}
	if len(rule.Include) == 0 {
		return fmt.Errorf("include is required")
	}
	if rule.Match.Type != "regex" {
		return fmt.Errorf("unsupported match type %q", rule.Match.Type)
	}
	if rule.Match.Pattern == "" {
		return fmt.Errorf("match pattern is required")
	}
	if _, err := regexp.Compile(rule.Match.Pattern); err != nil {
		return fmt.Errorf("invalid regex %q: %w", rule.Match.Pattern, err)
	}
	if rule.Message == "" {
		return fmt.Errorf("message is required")
	}
	return nil
}

func isYAML(name string) bool {
	lower := strings.ToLower(name)
	return strings.HasSuffix(lower, ".yaml") || strings.HasSuffix(lower, ".yml")
}

func yamlFiles(files []string) []string {
	var yaml []string
	for _, file := range files {
		if isYAML(file) {
			yaml = append(yaml, file)
		}
	}
	return yaml
}
