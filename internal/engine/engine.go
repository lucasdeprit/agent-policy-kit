package engine

import (
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/lucasdeprit/agent-policy-kit/internal/diff"
	"github.com/lucasdeprit/agent-policy-kit/internal/rules"
)

type Violation struct {
	RuleID     string
	File       string
	Line       int
	Text       string
	Message    string
	Expected   string
	Suggestion string
}

func Review(changes []diff.Change, ruleSet []rules.Rule) []Violation {
	compiled := compileRules(ruleSet)
	var violations []Violation

	for _, change := range changes {
		path := filepath.ToSlash(change.File)
		for _, rule := range compiled {
			if !rule.matchesPath(path) || !rule.pattern.MatchString(change.Text) {
				continue
			}

			violations = append(violations, Violation{
				RuleID:     rule.ID,
				File:       path,
				Line:       change.Line,
				Text:       change.Text,
				Message:    rule.Message,
				Expected:   rule.Expected,
				Suggestion: rule.Suggestion,
			})
		}
	}

	sort.SliceStable(violations, func(i, j int) bool {
		if violations[i].File != violations[j].File {
			return violations[i].File < violations[j].File
		}
		if violations[i].Line != violations[j].Line {
			return violations[i].Line < violations[j].Line
		}
		return violations[i].RuleID < violations[j].RuleID
	})

	return violations
}

type compiledRule struct {
	rules.Rule
	pattern *regexp.Regexp
}

func compileRules(ruleSet []rules.Rule) []compiledRule {
	compiled := make([]compiledRule, 0, len(ruleSet))
	for _, rule := range ruleSet {
		compiled = append(compiled, compiledRule{Rule: rule, pattern: regexp.MustCompile(rule.Match.Pattern)})
	}
	return compiled
}

func (r compiledRule) matchesPath(path string) bool {
	if !matchesAny(r.Include, path) {
		return false
	}
	return !matchesAny(r.Exclude, path)
}

func matchesAny(patterns []string, path string) bool {
	for _, pattern := range patterns {
		if globMatch(pattern, path) {
			return true
		}
	}
	return false
}

func globMatch(pattern, path string) bool {
	pattern = filepath.ToSlash(pattern)
	path = filepath.ToSlash(path)

	regex := globToRegex(pattern)
	return regexp.MustCompile(regex).MatchString(path)
}

func globToRegex(pattern string) string {
	var b strings.Builder
	b.WriteString("^")
	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '*':
			if i+1 < len(pattern) && pattern[i+1] == '*' {
				if i+2 < len(pattern) && pattern[i+2] == '/' {
					b.WriteString("(?:.*/)?")
					i += 2
				} else {
					b.WriteString(".*")
					i++
				}
			} else {
				b.WriteString("[^/]*")
			}
		case '?':
			b.WriteString("[^/]")
		default:
			b.WriteString(regexp.QuoteMeta(string(pattern[i])))
		}
	}
	b.WriteString("$")
	return b.String()
}
