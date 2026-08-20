package workspace

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

// ListDirectory lists the immediate children of a directory.
func ListDirectory(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in ListDirectoryInput,
) (*mcp.CallToolResult, ListDirectoryOutput, error) {
	if strings.TrimSpace(in.Path) == "" {
		return nil, ListDirectoryOutput{}, fmt.Errorf("path is required; provide a directory path")
	}

	cleanPath := filepath.Clean(in.Path)

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
