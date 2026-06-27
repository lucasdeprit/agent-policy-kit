package diff

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestParseUnifiedOnlyAddedLines(t *testing.T) {
	raw := []byte(`diff --git a/lib/page.dart b/lib/page.dart
index 1111111..2222222 100644
--- a/lib/page.dart
+++ b/lib/page.dart
@@ -10,2 +10,2 @@
-old line
+console.log("debug")
 context
@@ -20,0 +21,1 @@
+logger.info("ok")
`)

	changes, err := ParseUnified(raw)
	if err != nil {
		t.Fatalf("ParseUnified returned error: %v", err)
	}

	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}
	if changes[0].File != "lib/page.dart" || changes[0].Line != 10 || changes[0].Text != `console.log("debug")` {
		t.Fatalf("unexpected first change: %#v", changes[0])
	}
	if changes[1].File != "lib/page.dart" || changes[1].Line != 21 || changes[1].Text != `logger.info("ok")` {
		t.Fatalf("unexpected second change: %#v", changes[1])
	}
}

func TestWorkingTreeReadsTargetRepo(t *testing.T) {
	target := t.TempDir()
	runGit(t, target, "init")
	writeFile(t, filepath.Join(target, "tracked.txt"), "old\n")
	runGit(t, target, "add", ".")
	runGit(t, target, "-c", "user.name=Test", "-c", "user.email=test@example.com", "commit", "-m", "initial")
	writeFile(t, filepath.Join(target, "tracked.txt"), "old\nfoo\n")
	writeFile(t, filepath.Join(target, "new.txt"), "bar\n")

	changes, err := WorkingTree(target)
	if err != nil {
		t.Fatalf("WorkingTree returned error: %v", err)
	}
	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %#v", changes)
	}
	if changes[0].File != "tracked.txt" || changes[0].Text != "foo" {
		t.Fatalf("unexpected tracked change: %#v", changes[0])
	}
	if changes[1].File != "new.txt" || changes[1].Text != "bar" {
		t.Fatalf("unexpected untracked change: %#v", changes[1])
	}
}

func writeFile(t *testing.T, path, content string) {
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
