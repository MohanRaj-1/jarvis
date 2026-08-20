package workspace

import (
	"context"
	"fmt"
	"os"
	"strings"

	internalworkspace "jarvis/internal/workspace"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ListDirectoryInput contains the directory path to list.
type ListDirectoryInput struct {
	Path string `json:"path" jsonschema:"Absolute or relative directory path"`
}

// DirectoryEntry describes a directory child.
type DirectoryEntry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"is_dir"`
}

// ListDirectoryOutput contains the immediate children of a directory.
type ListDirectoryOutput struct {
	Entries []DirectoryEntry `json:"entries"`
}

// ListDirectoryTool lists directories in its workspace.
type ListDirectoryTool struct {
	workspace internalworkspace.Workspace
}

// NewListDirectoryTool creates a directory-listing tool limited to w.
func NewListDirectoryTool(w internalworkspace.Workspace) ListDirectoryTool {
	return ListDirectoryTool{workspace: w}
}

// ListDirectory lists the immediate children of a directory.
func (t ListDirectoryTool) Handle(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in ListDirectoryInput,
) (*mcp.CallToolResult, ListDirectoryOutput, error) {
	if strings.TrimSpace(in.Path) == "" {
		return nil, ListDirectoryOutput{}, fmt.Errorf("path is required; provide a directory path")
	}

	cleanPath, err := t.workspace.Resolve(in.Path)
	if err != nil {
		return nil, ListDirectoryOutput{}, err
	}

	info, err := os.Stat(cleanPath)
	if err != nil {
		return nil, ListDirectoryOutput{}, fmt.Errorf(
			"cannot access directory %q: %w",
			cleanPath,
			err,
		)
	}

	if !info.IsDir() {
		return nil, ListDirectoryOutput{}, fmt.Errorf(
			"%q is a file, not a directory; provide a directory path",
			cleanPath,
		)
	}

	dirEntries, err := os.ReadDir(cleanPath)
	if err != nil {
		return nil, ListDirectoryOutput{}, fmt.Errorf(
			"cannot read directory %q: %w",
			cleanPath,
			err,
		)
	}

	entries := make([]DirectoryEntry, 0, len(dirEntries))
	for _, entry := range dirEntries {
		entries = append(entries, DirectoryEntry{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
		})
	}

	return nil, ListDirectoryOutput{
		Entries: entries,
	}, nil
}
