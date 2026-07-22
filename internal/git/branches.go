package git

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	gitlib "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// BranchesResult describes the current branch and all local branches in a repository.
type BranchesResult struct {
	Current string
	Local   []string
}

// RepositoryBranches returns the current branch and all local branches in repoPath.
func RepositoryBranches(repoPath string) (*BranchesResult, error) {
	if strings.TrimSpace(repoPath) == "" {
		return nil, fmt.Errorf("repository path is required")
	}

	cleanPath := filepath.Clean(repoPath)
	info, err := os.Stat(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("access repository path %q: %w", cleanPath, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("repository path %q is not a directory", cleanPath)
	}

	repository, err := gitlib.PlainOpen(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("open Git repository %q: %w", cleanPath, err)
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
