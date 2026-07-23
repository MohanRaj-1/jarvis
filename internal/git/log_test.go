package git

import (
	"path/filepath"
	"testing"
	"time"

	gitlib "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestLog(t *testing.T) {
	repoPath := t.TempDir()
	repository, err := gitlib.PlainInit(repoPath, false)
	if err != nil {
		t.Fatalf("PlainInit(%q) returned an error: %v", repoPath, err)
	}
	worktree, err := repository.Worktree()
	if err != nil {
		t.Fatalf("Worktree() returned an error: %v", err)
	}

	commitTime := time.Date(2026, time.July, 23, 9, 15, 0, 0, time.UTC)
	for _, message := range []string{"first commit", "second commit"} {
		writeStatusTestFile(t, repoPath, "file.txt", message)
		if _, err := worktree.Add("file.txt"); err != nil {
			t.Fatalf("add file: %v", err)
		}
		if _, err := worktree.Commit(message, &gitlib.CommitOptions{
			Author: &object.Signature{Name: "Mohan Raj", Email: "mohan@example.com", When: commitTime},
		}); err != nil {
			t.Fatalf("commit file: %v", err)
		}
		commitTime = commitTime.Add(time.Minute)
	}

	commits, err := Log(repoPath, 1)
	if err != nil {
		t.Fatalf("Log(%q, 1) returned an error: %v", repoPath, err)
	}
	if len(commits) != 1 {
		t.Fatalf("Log(%q, 1) returned %d commits, want 1", repoPath, len(commits))
	}
	commit := commits[0]
	if commit.Message != "second commit" || commit.Author != "Mohan Raj" || commit.Email != "mohan@example.com" {
		t.Errorf("commit = %#v, want latest commit metadata", commit)
	}
	if !commit.Date.Equal(time.Date(2026, time.July, 23, 9, 16, 0, 0, time.UTC)) {
		t.Errorf("Date = %v, want latest commit time", commit.Date)
	}
	if len(commit.Hash) != 7 {
		t.Errorf("Hash = %q, want a 7-character abbreviated hash", commit.Hash)
	}
}

func TestLogUsesDefaultLimitAndRejectsNegativeLimit(t *testing.T) {
	repoPath := t.TempDir()
	if _, err := gitlib.PlainInit(repoPath, false); err != nil {
		t.Fatalf("PlainInit(%q) returned an error: %v", repoPath, err)
	}
	commits, err := Log(repoPath, 0)
	if err != nil {
		t.Fatalf("Log(%q, 0) returned an error: %v", repoPath, err)
	}
	if len(commits) != 0 {
		t.Errorf("Log(%q, 0) returned %d commits, want 0", repoPath, len(commits))
	}
	if _, err := Log(repoPath, -1); err == nil {
		t.Error("Log with a negative limit returned nil error")
	}
	if _, err := Log(filepath.Join(t.TempDir(), "not-a-repository"), 1); err == nil {
		t.Error("Log with an invalid repository path returned nil error")
	}
}
