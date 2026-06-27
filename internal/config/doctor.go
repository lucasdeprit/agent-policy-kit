package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lucasdeprit/agent-policy-kit/internal/parser"
)

type Check struct {
	ID      string
	Message string
	Path    string
	Fix     string
}

type DoctorReport struct {
	Options   Options
	Git       bool
	RuleFiles []string
	Failures  []Check
}

func Doctor(options Options) DoctorReport {
	report := DoctorReport{Options: options}

	if info, err := os.Stat(options.Target); err != nil {
		report.Failures = append(report.Failures, Check{ID: "D001", Message: "target path does not exist", Path: options.Target, Fix: "pass --target pointing to the project to validate"})
	} else if !info.IsDir() {
		report.Failures = append(report.Failures, Check{ID: "D002", Message: "target path is not a directory", Path: options.Target, Fix: "pass --target pointing to a directory"})
	}

	if len(options.RuleSources) == 0 {
		report.Failures = append(report.Failures, Check{ID: "D003", Message: "no rules sources configured", Fix: "pass --rules or create ./rules"})
	}

	for _, source := range options.RuleSources {
		files, err := parser.RuleFilesFromSource(source)
		if err != nil {
			report.Failures = append(report.Failures, Check{ID: "D004", Message: "rules source is not readable", Path: source, Fix: "pass --rules pointing to a readable rules directory, file, or glob"})
			continue
		}
		if len(files) == 0 {
			report.Failures = append(report.Failures, Check{ID: "D005", Message: "rules source has no yaml files", Path: source, Fix: "add .yaml/.yml rule files or pass another --rules path"})
			continue
		}
		report.RuleFiles = append(report.RuleFiles, files...)
	}

	if len(report.RuleFiles) > 0 {
		if _, err := parser.LoadRuleSources(options.RuleSources); err != nil {
			report.Failures = append(report.Failures, Check{ID: "D006", Message: fmt.Sprintf("rules failed to parse: %v", err), Fix: "fix invalid YAML or unsupported rule fields"})
		}
	}

	report.Git = PathExists(options.Target) && isGitWorkTree(options.Target)
	if PathExists(options.Target) && !report.Git {
		report.Failures = append(report.Failures, Check{ID: "D007", Message: "target is not a git work tree", Path: options.Target, Fix: "pass --target pointing inside a git repository"})
	}

	for _, source := range options.RuleSources {
		if samePath(options.Target, source) {
			report.Failures = append(report.Failures, Check{ID: "D008", Message: "target and rules path are the same", Path: source, Fix: "pass --target for the client project and --rules for the policy rules"})
		}
	}

	if looksLikePolicyKit(options.Target) {
		report.Failures = append(report.Failures, Check{ID: "D009", Message: "target looks like the Agent Policy Kit repository", Path: options.Target, Fix: "pass --target pointing to the project you want to validate"})
	}

	return report
}

func isGitWorkTree(target string) bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = target
	return cmd.Run() == nil
}

func samePath(a, b string) bool {
	absA, errA := filepath.Abs(a)
	absB, errB := filepath.Abs(b)
	if errA != nil || errB != nil {
		return false
	}
	return filepath.Clean(absA) == filepath.Clean(absB)
}

func looksLikePolicyKit(target string) bool {
	if !exists(target, "go.mod") || !exists(target, "cmd") || !exists(target, "internal") {
		return false
	}
	content, err := os.ReadFile(filepath.Join(target, "go.mod"))
	return err == nil && strings.Contains(string(content), "github.com/lucasdeprit/agent-policy-kit")
}

func exists(root, name string) bool {
	_, err := os.Stat(filepath.Join(root, name))
	return err == nil
}
