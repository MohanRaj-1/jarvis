// Package git provides Git repository helpers.
package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gitlib "github.com/go-git/go-git/v5"
)

// CurrentBranch returns the name of the branch currently checked out in repoPath.
func CurrentBranch(repoPath string) (string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return "", fmt.Errorf("repository path is required")
	}

	cleanPath := filepath.Clean(repoPath)
	info, err := os.Stat(cleanPath)
	if err != nil {
		return "", fmt.Errorf("access repository path %q: %w", cleanPath, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("repository path %q is not a directory", cleanPath)
	}

	repository, err := gitlib.PlainOpen(cleanPath)
	if err != nil {
		return "", fmt.Errorf("open Git repository %q: %w", cleanPath, err)
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
