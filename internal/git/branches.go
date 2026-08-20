package git

import (
	"fmt"
	"sort"

	"github.com/go-git/go-git/v5/plumbing"
)

// BranchesResult describes the current branch and all local branches in a repository.
type BranchesResult struct {
	Current string
	Local   []string
}

// RepositoryBranches returns the current branch and all local branches in repoPath.
func RepositoryBranches(repoPath string) (*BranchesResult, error) {
	repository, cleanPath, err := openRepository(repoPath)
	if err != nil {
		return nil, err
	}

	head, err := repository.Head()
	if err != nil {
		return nil, fmt.Errorf("read current branch for %q: %w", cleanPath, err)
	}
	if !head.Name().IsBranch() {
		return nil, fmt.Errorf("repository %q has a detached HEAD", cleanPath)
	}

	result := &BranchesResult{Current: head.Name().Short()}
	branches, err := repository.Branches()
	if err != nil {
		return nil, fmt.Errorf("list local branches for %q: %w", cleanPath, err)
	}
	if err := branches.ForEach(func(branch *plumbing.Reference) error {
		result.Local = append(result.Local, branch.Name().Short())
		return nil
	}); err != nil {
		return nil, fmt.Errorf("list local branches for %q: %w", cleanPath, err)
	}

	sort.Strings(result.Local)
	return result, nil
}
