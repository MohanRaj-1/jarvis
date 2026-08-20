package git

import (
	"path/filepath"
	"testing"
	"time"

	gitlib "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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

	commits, err := newTestRepository(t, repoPath).Log(repoPath, 1)
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
	repository := newTestRepository(t, repoPath)
	commits, err := repository.Log(repoPath, 0)
	if err != nil {
		t.Fatalf("Log(%q, 0) returned an error: %v", repoPath, err)
	}
	if len(commits) != 0 {
		t.Errorf("Log(%q, 0) returned %d commits, want 0", repoPath, len(commits))
	}
	if _, err := repository.Log(repoPath, -1); err == nil {
		t.Error("Log with a negative limit returned nil error")
	}
	if _, err := repository.Log(filepath.Join(t.TempDir(), "not-a-repository"), 1); err == nil {
		t.Error("Log with an invalid repository path returned nil error")
	}
}

func TestLogRange(t *testing.T) {
	repoPath := t.TempDir()
	repository, err := gitlib.PlainInit(repoPath, false)
	if err != nil {
		t.Fatalf("PlainInit(%q) returned an error: %v", repoPath, err)
	}
	worktree, err := repository.Worktree()
	if err != nil {
		t.Fatalf("Worktree() returned an error: %v", err)
	}

	commit := func(message string) string {
		writeStatusTestFile(t, repoPath, "file.txt", message)
		if _, err := worktree.Add("file.txt"); err != nil {
			t.Fatalf("add file: %v", err)
		}
		hash, err := worktree.Commit(message, &gitlib.CommitOptions{Author: &object.Signature{Name: "Mohan Raj", Email: "mohan@example.com", When: time.Now()}})
		if err != nil {
			t.Fatalf("commit file: %v", err)
		}
		return hash.String()
	}

	from := commit("B")
	commit("C")
	to := commit("D")

	repositoryService := newTestRepository(t, repoPath)
	commits, err := repositoryService.LogRange(repoPath, from, to)
	if err != nil {
		t.Fatalf("LogRange() error = %v", err)
	}
	if len(commits) != 2 || commits[0].Message != "D" || commits[1].Message != "C" {
		t.Errorf("LogRange(%q, %q) = %#v, want D then C (the B..D range)", from, to, commits)
	}
	if len(commits[0].ChangedFiles) != 1 || commits[0].ChangedFiles[0].Status != "Modified" {
		t.Errorf("LogRange().ChangedFiles = %#v, want one modified file", commits[0].ChangedFiles)
	}
	if _, err := repositoryService.LogRange(repoPath, "missing", "HEAD"); err == nil {
		t.Error("LogRange() with an unknown start revision returned nil error")
	}
}

func TestLogRangeTagToHEAD(t *testing.T) {
	repoPath := t.TempDir()
	repository, err := gitlib.PlainInit(repoPath, false)
	if err != nil {
		t.Fatalf("PlainInit(%q) returned an error: %v", repoPath, err)
	}
	worktree, err := repository.Worktree()
	if err != nil {
		t.Fatalf("Worktree() returned an error: %v", err)
	}

	commit := func(message string) plumbing.Hash {
		writeStatusTestFile(t, repoPath, "file.txt", message)
		if _, err := worktree.Add("file.txt"); err != nil {
			t.Fatalf("add file: %v", err)
		}
		hash, err := worktree.Commit(message, &gitlib.CommitOptions{
			Author: &object.Signature{Name: "Mohan Raj", Email: "mohan@example.com", When: time.Now()},
		})
		if err != nil {
			t.Fatalf("commit file: %v", err)
		}
		return hash
	}

	taggedCommit := commit("release v0.5.2")
	if _, err := repository.CreateTag("v0.5.2", taggedCommit, nil); err != nil {
		t.Fatalf("CreateTag(v0.5.2) returned an error: %v", err)
	}
	commit("feature after release")
	commit("fix after release")

	commits, err := newTestRepository(t, repoPath).LogRange(repoPath, "v0.5.2", "HEAD")
	if err != nil {
		t.Fatalf("LogRange(v0.5.2, HEAD) returned an error: %v", err)
	}
	if len(commits) != 2 || commits[0].Message != "fix after release" || commits[1].Message != "feature after release" {
		t.Errorf("LogRange(v0.5.2, HEAD) = %#v, want commits after the v0.5.2 tag in reverse chronological order", commits)
	}
}

func TestLogRangeSameEndpoint(t *testing.T) {
	repoPath := t.TempDir()
	repository, err := gitlib.PlainInit(repoPath, false)
	if err != nil {
		t.Fatalf("PlainInit(%q) returned an error: %v", repoPath, err)
	}
	worktree, err := repository.Worktree()
	if err != nil {
		t.Fatalf("Worktree() returned an error: %v", err)
	}

	writeStatusTestFile(t, repoPath, "file.txt", "first commit")
	if _, err := worktree.Add("file.txt"); err != nil {
		t.Fatalf("add file: %v", err)
	}
	hash, err := worktree.Commit("first commit", &gitlib.CommitOptions{
		Author: &object.Signature{Name: "Mohan Raj", Email: "mohan@example.com", When: time.Now()},
	})
	if err != nil {
		t.Fatalf("commit file: %v", err)
	}

	commits, err := newTestRepository(t, repoPath).LogRange(repoPath, hash.String(), hash.String())
	if err != nil {
		t.Fatalf("LogRange() error = %v", err)
	}
	if len(commits) != 0 {
		t.Errorf("LogRange() = %#v, want no commits when both endpoints are the same", commits)
	}
}
