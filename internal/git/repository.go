package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gitlib "github.com/go-git/go-git/v5"
)

// openRepository validates repoPath and opens the Git repository it contains.
// It returns the cleaned path so callers can include it in operation-specific errors.
func openRepository(repoPath string) (*gitlib.Repository, string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return nil, "", fmt.Errorf("repository path is required")
	}

	cleanPath := filepath.Clean(repoPath)
	info, err := os.Stat(cleanPath)
	if err != nil {
		return nil, "", fmt.Errorf("access repository path %q: %w", cleanPath, err)
	}
	if !info.IsDir() {
		return nil, "", fmt.Errorf("repository path %q is not a directory", cleanPath)
	}

	repository, err := gitlib.PlainOpen(cleanPath)
	if err != nil {
		return nil, "", fmt.Errorf("open Git repository %q: %w", cleanPath, err)
	}

	return repository, cleanPath, nil
}
