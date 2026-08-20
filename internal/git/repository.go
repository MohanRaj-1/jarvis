package git

import (
	"fmt"
	"os"
	"strings"

	internalworkspace "jarvis/internal/workspace"

	gitlib "github.com/go-git/go-git/v5"
)

// Repository provides Git operations used by higher-level services.
// Implementations can be replaced in tests without depending on a real Git
// repository.
type Repository interface {
	Diff(repoPath string) (string, error)
	Show(repoPath, hash string) (*CommitDetails, error)
	LogRange(repoPath, from, to string) ([]ReleaseCommit, error)
}

// DefaultRepository is the production Repository implementation.
type DefaultRepository struct {
	workspace internalworkspace.Workspace
}

// NewDefaultRepository creates a Git repository implementation limited to workspace.
func NewDefaultRepository(workspace internalworkspace.Workspace) DefaultRepository {
	return DefaultRepository{workspace: workspace}
}

// openRepository validates repoPath and opens the Git repository it contains.
// It returns the cleaned path so callers can include it in operation-specific errors.
func (r DefaultRepository) openRepository(repoPath string) (*gitlib.Repository, string, error) {
	if strings.TrimSpace(repoPath) == "" {
		return nil, "", fmt.Errorf("repository path is required")
	}

	cleanPath, err := r.workspace.Resolve(repoPath)
	if err != nil {
		return nil, "", err
	}
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
