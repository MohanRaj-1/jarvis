package git

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/go-git/go-git/v5/plumbing"
)

// Diff returns a unified diff between the repository's current HEAD and
// working tree.
func Diff(repoPath string) (string, error) {
	repository, cleanPath, err := openRepository(repoPath)
	if err != nil {
		return "", err
	}

	if _, err := repository.Worktree(); err != nil {
		return "", fmt.Errorf("open worktree for %q: %w", cleanPath, err)
	}

	head, err := repository.Head()
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			return "", nil
		}
		return "", fmt.Errorf("read HEAD for %q: %w", cleanPath, err)
	}

	command := exec.Command("git", "-C", cleanPath, "diff", "--no-ext-diff", head.Hash().String(), "--")
	output, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("get diff against HEAD for %q: %w", cleanPath, err)
	}

	return string(output), nil
}
