// Package git provides Git repository helpers.
package git

import (
	"fmt"
)

// CurrentBranch returns the name of the branch currently checked out in repoPath.
func CurrentBranch(repoPath string) (string, error) {
	repository, cleanPath, err := openRepository(repoPath)
	if err != nil {
		return "", err
	}

	head, err := repository.Head()
	if err != nil {
		return "", fmt.Errorf("read current branch for %q: %w", cleanPath, err)
	}
	if !head.Name().IsBranch() {
		return "", fmt.Errorf("repository %q has a detached HEAD", cleanPath)
	}

	return head.Name().Short(), nil
}
