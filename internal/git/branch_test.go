package git

import (
	"path/filepath"
	"testing"

	gitlib "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func TestCurrentBranch(t *testing.T) {
	repoPath := t.TempDir()
	repository, err := gitlib.PlainInit(repoPath, false)
	if err != nil {
		t.Fatalf("PlainInit(%q) returned an error: %v", repoPath, err)
	}

	branchName := plumbing.NewBranchReferenceName("feature/git-intelligence")
	if err := repository.Storer.SetReference(plumbing.NewHashReference(branchName, plumbing.ZeroHash)); err != nil {
		t.Fatalf("set branch reference: %v", err)
	}
	if err := repository.Storer.SetReference(plumbing.NewSymbolicReference(plumbing.HEAD, branchName)); err != nil {
		t.Fatalf("set HEAD reference: %v", err)
	}

	branch, err := CurrentBranch(repoPath)
	if err != nil {
		t.Fatalf("CurrentBranch(%q) returned an error: %v", repoPath, err)
	}

	if branch != "feature/git-intelligence" {
		t.Errorf("CurrentBranch(%q) = %q, want %q", repoPath, branch, "feature/git-intelligence")
	}
}

func TestCurrentBranchRejectsInvalidRepositoryPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "not-a-repository")

	if _, err := CurrentBranch(path); err == nil {
		t.Errorf("CurrentBranch(%q) returned nil error", path)
	}
}
