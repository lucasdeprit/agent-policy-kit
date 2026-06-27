package diff

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Change struct {
	File string
	Line int
	Text string
}

func WorkingTree(target string) ([]Change, error) {
	raw, err := gitDiffHead(target)
	if err != nil {
		return nil, err
	}

	changes, err := ParseUnified(raw)
	if err != nil {
		return nil, err
	}

	untracked, err := untrackedChanges(target)
	if err != nil {
		return nil, err
	}
	changes = append(changes, untracked...)
	return changes, nil
}

func gitDiffHead(target string) ([]byte, error) {
	cmd := exec.Command("git", "diff", "--no-ext-diff", "--unified=0", "HEAD", "--")
	cmd.Dir = target
	out, err := cmd.Output()
	if err == nil {
		return out, nil
	}

	cmd = exec.Command("git", "diff", "--no-ext-diff", "--unified=0", "--")
	cmd.Dir = target
	out, fallbackErr := cmd.Output()
	if fallbackErr != nil {
		return nil, fmt.Errorf("could not read git diff: %w", err)
	}
	return out, nil
}

var hunkHeader = regexp.MustCompile(`^@@ -\d+(?:,\d+)? \+(\d+)(?:,\d+)? @@`)

func ParseUnified(raw []byte) ([]Change, error) {
	var changes []Change
	var file string
	newLine := 0

	scanner := bufio.NewScanner(bytes.NewReader(raw))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "+++ b/"):
			file = strings.TrimPrefix(line, "+++ b/")
		case strings.HasPrefix(line, "@@"):
			match := hunkHeader.FindStringSubmatch(line)
			if len(match) != 2 {
				return nil, fmt.Errorf("invalid hunk header: %s", line)
			}
			parsed, err := strconv.Atoi(match[1])
			if err != nil {
				return nil, fmt.Errorf("invalid hunk line number %q: %w", match[1], err)
			}
			newLine = parsed
		case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
			if file == "" || newLine == 0 {
				continue
			}
			changes = append(changes, Change{File: file, Line: newLine, Text: strings.TrimPrefix(line, "+")})
			newLine++
		case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---"):
			// Removed lines do not advance the new-file line counter.
		default:
			if newLine > 0 {
				newLine++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan diff: %w", err)
	}
	return changes, nil
}

func untrackedChanges(target string) ([]Change, error) {
	cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
	cmd.Dir = target
	out, err := cmd.Output()
	if err != nil {
		return nil, nil
	}

	var changes []Change
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		path := scanner.Text()
		absolutePath := filepath.Join(target, filepath.FromSlash(path))
		info, err := os.Stat(absolutePath)
		if err != nil || info.IsDir() {
			continue
		}

		content, err := os.ReadFile(absolutePath)
		if err != nil || bytes.IndexByte(content, 0) >= 0 {
			continue
		}

		lineScanner := bufio.NewScanner(bytes.NewReader(content))
		lineNo := 1
		for lineScanner.Scan() {
			changes = append(changes, Change{File: path, Line: lineNo, Text: lineScanner.Text()})
			lineNo++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan untracked files: %w", err)
	}
	return changes, nil
}
