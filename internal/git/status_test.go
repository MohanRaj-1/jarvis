package git

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	gitlib "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestRepositoryStatus(t *testing.T) {
	repoPath := t.TempDir()
	repository, err := gitlib.PlainInit(repoPath, false)
	if err != nil {
		t.Fatalf("PlainInit(%q) returned an error: %v", repoPath, err)
	}

	writeStatusTestFile(t, repoPath, "tracked.txt", "initial\n")
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

	writeStatusTestFile(t, repoPath, "tracked.txt", "modified\n")
	writeStatusTestFile(t, repoPath, "staged.txt", "staged\n")
	if _, err := worktree.Add("staged.txt"); err != nil {
		t.Fatalf("add staged file: %v", err)
	}
	writeStatusTestFile(t, repoPath, "untracked.txt", "untracked\n")

	status, err := RepositoryStatus(repoPath)
	if err != nil {
		t.Fatalf("RepositoryStatus(%q) returned an error: %v", repoPath, err)
	}
	if status.Branch != "master" {
		t.Errorf("Branch = %q, want %q", status.Branch, "master")
	}
	if want := []string{"tracked.txt"}; !reflect.DeepEqual(status.Modified, want) {
		t.Errorf("Modified = %#v, want %#v", status.Modified, want)
	}
	if want := []string{"staged.txt"}; !reflect.DeepEqual(status.Staged, want) {
		t.Errorf("Staged = %#v, want %#v", status.Staged, want)
	}
	if want := []string{"untracked.txt"}; !reflect.DeepEqual(status.Untracked, want) {
		t.Errorf("Untracked = %#v, want %#v", status.Untracked, want)
	}
}

func TestRepositoryStatusRejectsInvalidRepositoryPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "not-a-repository")

	if _, err := RepositoryStatus(path); err == nil {
		t.Errorf("RepositoryStatus(%q) returned nil error", path)
	}
}

func writeStatusTestFile(t *testing.T, repoPath, name, contents string) {
	t.Helper()
	path := filepath.Join(repoPath, name)
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write %q: %v", path, err)
	}
}
