package git

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	gitlib "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestDiff(t *testing.T) {
	repoPath := t.TempDir()
	repository, err := gitlib.PlainInit(repoPath, false)
	if err != nil {
		t.Fatalf("PlainInit(%q) returned an error: %v", repoPath, err)
	}

	writeDiffTestFile(t, repoPath, "tracked.txt", "before\n")
	worktree, err := repository.Worktree()
	if err != nil {
		t.Fatalf("Worktree() returned an error: %v", err)
	}
	if _, err := worktree.Add("tracked.txt"); err != nil {
		t.Fatalf("add tracked file: %v", err)
	}
	if _, err := worktree.Commit("initial commit", &gitlib.CommitOptions{
		Author: &object.Signature{Name: "Test User", Email: "test@example.com", When: time.Now()},
	}); err != nil {
		t.Fatalf("commit tracked file: %v", err)
	}

	writeDiffTestFile(t, repoPath, "tracked.txt", "after\n")
	writeDiffTestFile(t, repoPath, "staged.txt", "staged\n")
	if _, err := worktree.Add("staged.txt"); err != nil {
		t.Fatalf("add staged file: %v", err)
	}
	diff, err := Diff(repoPath)
	if err != nil {
		t.Fatalf("Diff(%q) returned an error: %v", repoPath, err)
	}

	for _, want := range []string{
		"diff --git a/tracked.txt b/tracked.txt",
		"--- a/tracked.txt",
		"+++ b/tracked.txt",
		"-before",
		"+after",
		"diff --git a/staged.txt b/staged.txt",
		"+staged",
	} {
		if !strings.Contains(diff, want) {
			t.Errorf("Diff(%q) = %q, want it to contain %q", repoPath, diff, want)
		}
	}
}

func TestDiffCleanRepository(t *testing.T) {
	repoPath := t.TempDir()
	if _, err := gitlib.PlainInit(repoPath, false); err != nil {
		t.Fatalf("PlainInit(%q) returned an error: %v", repoPath, err)
	}

	diff, err := Diff(repoPath)
	if err != nil {
		t.Fatalf("Diff(%q) returned an error: %v", repoPath, err)
	}
	if diff != "" {
		t.Errorf("Diff(%q) = %q, want an empty diff", repoPath, diff)
	}
}

func TestDiffRejectsInvalidRepository(t *testing.T) {
	path := filepath.Join(t.TempDir(), "not-a-repository")

	if _, err := Diff(path); err == nil {
		t.Errorf("Diff(%q) returned nil error", path)
	}
}

func TestDiffRejectsNonGitDirectory(t *testing.T) {
	path := t.TempDir()

	if _, err := Diff(path); err == nil {
		t.Errorf("Diff(%q) returned nil error", path)
	}
}

func writeDiffTestFile(t *testing.T, repoPath, name, contents string) {
	t.Helper()
	path := filepath.Join(repoPath, name)
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write %q: %v", path, err)
	}
}
