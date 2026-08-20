// Package workspace restricts filesystem operations to a configured workspace.
package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Workspace is the canonical absolute directory that filesystem operations may access.
type Workspace struct {
	root string
}

// New creates a workspace rooted at root. The root must exist and be a directory.
func New(root string) (Workspace, error) {
	if strings.TrimSpace(root) == "" {
		return Workspace{}, errors.New("workspace root is required")
	}

	absRoot, err := filepath.Abs(filepath.Clean(root))
	if err != nil {
		return Workspace{}, fmt.Errorf("resolve workspace root: %w", err)
	}

	resolvedRoot, err := filepath.EvalSymlinks(absRoot)
	if err != nil {
		return Workspace{}, fmt.Errorf("resolve workspace root symlinks: %w", err)
	}

	info, err := os.Stat(resolvedRoot)
	if err != nil {
		return Workspace{}, fmt.Errorf("stat workspace root: %w", err)
	}
	if !info.IsDir() {
		return Workspace{}, fmt.Errorf("workspace root %q is not a directory", root)
	}

	return Workspace{root: resolvedRoot}, nil
}

// Resolve returns path as a canonical absolute path when it is contained by the workspace.
// Relative paths are resolved from the workspace root. Existing symlinks are resolved before
// the containment check, so a symlink cannot be used to escape the workspace.
func (w Workspace) Resolve(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", errors.New("path is required")
	}
	if w.root == "" {
		return "", errors.New("workspace root is required")
	}

	candidate := path
	if !filepath.IsAbs(candidate) {
		candidate = filepath.Join(w.root, candidate)
	}

	absPath, err := filepath.Abs(filepath.Clean(candidate))
	if err != nil {
		return "", fmt.Errorf("resolve path: %w", err)
	}
	resolvedPath, err := evalSymlinksWherePossible(absPath)
	if err != nil {
		return "", fmt.Errorf("resolve path symlinks: %w", err)
	}
	if !within(w.root, resolvedPath) {
		return "", fmt.Errorf("path %q resolves outside workspace", path)
	}

	return resolvedPath, nil
}

func within(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// evalSymlinksWherePossible resolves the longest existing path prefix. This permits callers
// to resolve a not-yet-created leaf while still checking every existing symlink in its path.
func evalSymlinksWherePossible(path string) (string, error) {
	var missing []string
	current := path
	for {
		resolved, err := filepath.EvalSymlinks(current)
		if err == nil {
			for i := len(missing) - 1; i >= 0; i-- {
				resolved = filepath.Join(resolved, missing[i])
			}
			return filepath.Clean(resolved), nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return "", err
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", err
		}
		missing = append(missing, filepath.Base(current))
		current = parent
	}
}
