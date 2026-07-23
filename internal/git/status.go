package git

import (
	"fmt"
	"sort"

	gitlib "github.com/go-git/go-git/v5"
)

// Status describes the changes currently present in a Git repository.
type Status struct {
	Branch    string
	Modified  []string
	Staged    []string
	Untracked []string
}

// RepositoryStatus returns the current branch and changed files in repoPath.
func RepositoryStatus(repoPath string) (*Status, error) {
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

	worktree, err := repository.Worktree()
	if err != nil {
		return nil, fmt.Errorf("open worktree for %q: %w", cleanPath, err)
	}
	gitStatus, err := worktree.Status()
	if err != nil {
		return nil, fmt.Errorf("get status for %q: %w", cleanPath, err)
	}

	status := &Status{Branch: head.Name().Short()}
	for path, fileStatus := range gitStatus {
		switch {
		case fileStatus.Worktree == gitlib.Untracked:
			status.Untracked = append(status.Untracked, path)
		case fileStatus.Staging != gitlib.Unmodified:
			status.Staged = append(status.Staged, path)
		}
		if fileStatus.Worktree != gitlib.Unmodified && fileStatus.Worktree != gitlib.Untracked {
			status.Modified = append(status.Modified, path)
		}
	}

	sort.Strings(status.Modified)
	sort.Strings(status.Staged)
	sort.Strings(status.Untracked)
	return status, nil
}
