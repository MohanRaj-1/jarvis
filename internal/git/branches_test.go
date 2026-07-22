package git

import (
	"reflect"
	"testing"

	gitlib "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func TestRepositoryBranches(t *testing.T) {
	repoPath := t.TempDir()
	repository, err := gitlib.PlainInit(repoPath, false)
	if err != nil {
		t.Fatalf("PlainInit(%q) returned an error: %v", repoPath, err)
	}

	for _, name := range []string{"develop", "feature/git-intelligence"} {
		branch := plumbing.NewBranchReferenceName(name)
		if err := repository.Storer.SetReference(plumbing.NewHashReference(branch, plumbing.ZeroHash)); err != nil {
			t.Fatalf("set branch reference %q: %v", name, err)
		}
	}
	current := plumbing.NewBranchReferenceName("feature/git-intelligence")
	if err := repository.Storer.SetReference(plumbing.NewSymbolicReference(plumbing.HEAD, current)); err != nil {
		t.Fatalf("set HEAD reference: %v", err)
	}

	result, err := RepositoryBranches(repoPath)
	if err != nil {
		t.Fatalf("RepositoryBranches(%q) returned an error: %v", repoPath, err)
	}
	if result.Current != "feature/git-intelligence" {
		t.Errorf("Current = %q, want %q", result.Current, "feature/git-intelligence")
	}
	if want := []string{"develop", "feature/git-intelligence", "master"}; !reflect.DeepEqual(result.Local, want) {
		t.Errorf("Local = %#v, want %#v", result.Local, want)
	}
}
