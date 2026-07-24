package git

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	gitlib "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestShow(t *testing.T) {
	repoPath := t.TempDir()
	repository, err := gitlib.PlainInit(repoPath, false)
	if err != nil {
		t.Fatalf("PlainInit(%q): %v", repoPath, err)
	}
	worktree, err := repository.Worktree()
	if err != nil {
		t.Fatalf("Worktree(): %v", err)
	}

	for name, contents := range map[string]string{
		"modified.txt": "before", "deleted.txt": "delete me", "old-name.txt": "rename me",
	} {
		writeStatusTestFile(t, repoPath, name, contents)
		if _, err := worktree.Add(name); err != nil {
			t.Fatalf("add %q: %v", name, err)
		}
	}
	if _, err := worktree.Commit("initial commit", showCommitOptions(time.Date(2026, time.July, 24, 10, 0, 0, 0, time.UTC))); err != nil {
		t.Fatalf("create initial commit: %v", err)
	}

	writeStatusTestFile(t, repoPath, "modified.txt", "after")
	writeStatusTestFile(t, repoPath, "added.txt", "new file")
	if err := os.Remove(filepath.Join(repoPath, "deleted.txt")); err != nil {
		t.Fatalf("remove deleted file: %v", err)
	}
	if err := os.Rename(filepath.Join(repoPath, "old-name.txt"), filepath.Join(repoPath, "new-name.txt")); err != nil {
		t.Fatalf("rename file: %v", err)
	}
	for _, name := range []string{"modified.txt", "added.txt", "new-name.txt"} {
		if _, err := worktree.Add(name); err != nil {
			t.Fatalf("add %q: %v", name, err)
		}
	}
	for _, name := range []string{"deleted.txt", "old-name.txt"} {
		if _, err := worktree.Remove(name); err != nil {
			t.Fatalf("remove %q from index: %v", name, err)
		}
	}
	hash, err := worktree.Commit("update files", showCommitOptions(time.Date(2026, time.July, 24, 10, 1, 0, 0, time.UTC)))
	if err != nil {
		t.Fatalf("create update commit: %v", err)
	}

	details, err := Show(repoPath, hash.String()[:7])
	if err != nil {
		t.Fatalf("Show(%q, %q): %v", repoPath, hash.String()[:7], err)
	}
	if details.Hash != hash.String()[:7] || details.Author != "Mohan Raj" || details.Email != "mohan@example.com" || details.Message != "update files" {
		t.Errorf("commit details = %#v, want update commit metadata", details.Commit)
	}
	if !details.Date.Equal(time.Date(2026, time.July, 24, 10, 1, 0, 0, time.UTC)) {
		t.Errorf("Date = %v, want commit time", details.Date)
	}
	if len(details.Parents) != 1 {
		t.Fatalf("Parents = %#v, want one parent", details.Parents)
	}

	statuses := make(map[string]string, len(details.Files))
	for _, file := range details.Files {
		statuses[file.Path] = file.Status
	}
	for path, want := range map[string]string{
		"added.txt": "Added", "modified.txt": "Modified", "deleted.txt": "Deleted", "new-name.txt": "Renamed",
	} {
		if got := statuses[path]; got != want {
			t.Errorf("status for %q = %q, want %q (all files: %#v)", path, got, want, details.Files)
		}
	}
}

func TestShowRejectsInvalidOrEmptyHashAndInvalidRepository(t *testing.T) {
	repoPath := t.TempDir()
	if _, err := gitlib.PlainInit(repoPath, false); err != nil {
		t.Fatalf("PlainInit(%q): %v", repoPath, err)
	}
	for _, hash := range []string{"", "not-a-commit"} {
		if _, err := Show(repoPath, hash); err == nil {
			t.Errorf("Show(%q, %q) returned nil error", repoPath, hash)
		}
	}
	if _, err := Show(filepath.Join(t.TempDir(), "not-a-repository"), "abcdef1"); err == nil {
		t.Error("Show with an invalid repository returned nil error")
	}
}

func TestShowReturnsAllMergeParents(t *testing.T) {
	repoPath := t.TempDir()
	repository, err := gitlib.PlainInit(repoPath, false)
	if err != nil {
		t.Fatalf("PlainInit(%q): %v", repoPath, err)
	}
	worktree, err := repository.Worktree()
	if err != nil {
		t.Fatalf("Worktree(): %v", err)
	}
	writeStatusTestFile(t, repoPath, "file.txt", "base")
	if _, err := worktree.Add("file.txt"); err != nil {
		t.Fatalf("add base file: %v", err)
	}
	base, err := worktree.Commit("base", showCommitOptions(time.Date(2026, time.July, 24, 11, 0, 0, 0, time.UTC)))
	if err != nil {
		t.Fatalf("create base commit: %v", err)
	}
	writeStatusTestFile(t, repoPath, "file.txt", "first parent")
	if _, err := worktree.Add("file.txt"); err != nil {
		t.Fatalf("add first parent change: %v", err)
	}
	firstParent, err := worktree.Commit("first parent", showCommitOptions(time.Date(2026, time.July, 24, 11, 1, 0, 0, time.UTC)))
	if err != nil {
		t.Fatalf("create first parent: %v", err)
	}
	mergeOptions := showCommitOptions(time.Date(2026, time.July, 24, 11, 2, 0, 0, time.UTC))
	mergeOptions.Parents = []plumbing.Hash{firstParent, base}
	mergeOptions.AllowEmptyCommits = true
	merge, err := worktree.Commit("merge commit", mergeOptions)
	if err != nil {
		t.Fatalf("create merge commit: %v", err)
	}

	details, err := Show(repoPath, merge.String()[:7])
	if err != nil {
		t.Fatalf("Show merge commit: %v", err)
	}
	if len(details.Parents) != 2 {
		t.Fatalf("Parents = %#v, want two parents", details.Parents)
	}
	if details.Parents[0] != firstParent.String()[:7] || details.Parents[1] != base.String()[:7] {
		t.Errorf("Parents = %#v, want [%q %q]", details.Parents, firstParent.String()[:7], base.String()[:7])
	}
}

func showCommitOptions(when time.Time) *gitlib.CommitOptions {
	return &gitlib.CommitOptions{Author: &object.Signature{
		Name: "Mohan Raj", Email: "mohan@example.com", When: when,
	}}
}
