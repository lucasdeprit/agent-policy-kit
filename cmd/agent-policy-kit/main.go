package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/lucasdeprit/agent-policy-kit/internal/config"
	"github.com/lucasdeprit/agent-policy-kit/internal/diff"
	"github.com/lucasdeprit/agent-policy-kit/internal/engine"
	"github.com/lucasdeprit/agent-policy-kit/internal/parser"
	"github.com/lucasdeprit/agent-policy-kit/internal/reporter"
)

//go:embed templates/rules
var templates embed.FS

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: agent-policy-kit <init|review|doctor>")
		return 1
	}

	switch args[0] {
	case "init":
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "usage: agent-policy-kit init")
			return 1
		}
		return runInit(".")
	case "review":
		return runReview(args[1:])
	case "doctor":
		return runDoctor(args[1:])
	default:
		fmt.Fprintln(os.Stderr, "usage: agent-policy-kit <init|review|doctor>")
		return 1
	}
}

func runReview(args []string) int {
	options, err := parseOptions("review", args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	rules, err := parser.LoadRuleSources(options.RuleSources)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	changes, err := diff.WorkingTree(options.Target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	changes = excludeRuleSourceChanges(changes, options.Target, options.RuleSources)

	violations := engine.Review(changes, rules)
	reporter.Text(os.Stdout, violations)
	if len(violations) > 0 {
		return 1
	}
	return 0
}

func excludeRuleSourceChanges(changes []diff.Change, target string, ruleSources []string) []diff.Change {
	exclusions := ruleSourceExclusions(ruleSources)
	if len(exclusions) == 0 {
		return changes
	}

	filtered := make([]diff.Change, 0, len(changes))
	for _, change := range changes {
		changePath := filepath.Join(target, filepath.FromSlash(change.File))
		if isExcludedPath(changePath, exclusions) {
			continue
		}
		filtered = append(filtered, change)
	}
	return filtered
}

type ruleSourceExclusion struct {
	Path  string
	IsDir bool
}

func ruleSourceExclusions(ruleSources []string) []ruleSourceExclusion {
	var exclusions []ruleSourceExclusion
	for _, source := range ruleSources {
		if strings.ContainsAny(source, "*?[") {
			matches, err := filepath.Glob(source)
			if err != nil {
				continue
			}
			for _, match := range matches {
				if exclusion, ok := newRuleSourceExclusion(match); ok {
					exclusions = append(exclusions, exclusion)
				}
			}
			continue
		}

		if exclusion, ok := newRuleSourceExclusion(source); ok {
			exclusions = append(exclusions, exclusion)
		}
	}
	return exclusions
}

func newRuleSourceExclusion(path string) (ruleSourceExclusion, bool) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return ruleSourceExclusion{}, false
	}
	info, err := os.Stat(absolutePath)
	if err != nil {
		return ruleSourceExclusion{}, false
	}
	return ruleSourceExclusion{Path: filepath.Clean(absolutePath), IsDir: info.IsDir()}, true
}

func isExcludedPath(path string, exclusions []ruleSourceExclusion) bool {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	absolutePath = filepath.Clean(absolutePath)

	for _, exclusion := range exclusions {
		if exclusion.IsDir {
			rel, err := filepath.Rel(exclusion.Path, absolutePath)
			if err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
				return true
			}
			continue
		}
		if absolutePath == exclusion.Path {
			return true
		}
	}
	return false
}

func runDoctor(args []string) int {
	options, err := parseOptions("doctor", args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	report := config.Doctor(options)
	printDoctor(os.Stdout, report)
	if len(report.Failures) > 0 {
		return 1
	}
	return 0
}

type ruleSources []string

func (s *ruleSources) String() string {
	return strings.Join(*s, ",")
}

func (s *ruleSources) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func parseOptions(command string, args []string) (config.Options, error) {
	flags := flag.NewFlagSet(command, flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	target := flags.String("target", ".", "project path whose diff is reviewed")
	var rules ruleSources
	flags.Var(&rules, "rules", "rules directory, file, or glob; can be repeated")
	if err := flags.Parse(args); err != nil {
		return config.Options{}, err
	}
	if flags.NArg() != 0 {
		return config.Options{}, fmt.Errorf("unexpected arguments: %s", strings.Join(flags.Args(), " "))
	}
	return config.Resolve(*target, rules)
}

func printDoctor(w io.Writer, report config.DoctorReport) {
	if len(report.Failures) == 0 {
		fmt.Fprintln(w, "OK doctor")
		fmt.Fprintln(w)
	} else {
		fmt.Fprintf(w, "FAILED doctor %d\n", len(report.Failures))
		fmt.Fprintln(w)
	}

	fmt.Fprintf(w, "target: %s\n", report.Options.Target)
	fmt.Fprintf(w, "rules: %s\n", strings.Join(report.Options.RuleSources, ", "))
	fmt.Fprintf(w, "rules_found: %d\n", len(report.RuleFiles))
	fmt.Fprintf(w, "git: %t\n", report.Git)

	for _, failure := range report.Failures {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "%s %s\n", failure.ID, failure.Message)
		if failure.Path != "" {
			fmt.Fprintf(w, "path: %s\n", failure.Path)
		}
		if failure.Fix != "" {
			fmt.Fprintf(w, "fix: %s\n", failure.Fix)
		}
	}
}

func runInit(root string) int {
	created, skipped, err := initProject(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Fprintf(os.Stdout, "initialized rules: created=%d skipped=%d\n", created, skipped)
	fmt.Fprintln(os.Stdout, "next: edit rules/*.yaml, then run agent-policy-kit review")
	return 0
}

func initProject(root string) (int, int, error) {
	created := 0
	skipped := 0

	err := fs.WalkDir(templates, "templates/rules", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel("templates", path)
		if err != nil {
			return err
		}
		target := filepath.Join(root, filepath.FromSlash(rel))

		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		content, err := templates.ReadFile(path)
		if err != nil {
			return err
		}

		_, statErr := os.Stat(target)
		if statErr == nil {
			skipped++
			return nil
		}
		if !os.IsNotExist(statErr) {
			return statErr
		}

		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(target, normalizeNewline(content), 0o644); err != nil {
			return err
		}
		created++
		return nil
	})

	return created, skipped, err
}

func normalizeNewline(content []byte) []byte {
	if strings.HasSuffix(string(content), "\n") {
		return content
	}
	return append(content, '\n')
}
