package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Options struct {
	Target      string
	RuleSources []string
}

func Resolve(target string, ruleSources []string) (Options, error) {
	if target == "" {
		target = "."
	}
	if len(ruleSources) == 0 {
		ruleSources = []string{"rules"}
	}

	absTarget, err := filepath.Abs(target)
	if err != nil {
		return Options{}, fmt.Errorf("resolve target %q: %w", target, err)
	}

	absRules := make([]string, 0, len(ruleSources))
	for _, source := range ruleSources {
		absSource, err := filepath.Abs(source)
		if err != nil {
			return Options{}, fmt.Errorf("resolve rules %q: %w", source, err)
		}
		absRules = append(absRules, absSource)
	}

	return Options{Target: absTarget, RuleSources: absRules}, nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
