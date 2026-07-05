package workspace

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ListDirectoryInput struct {
	Path string `json:"path" jsonschema:"Absolute or relative directory path"`
}

type DirectoryEntry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"is_dir"`
}

type ListDirectoryOutput struct {
	Entries []DirectoryEntry `json:"entries"`
}

func ListDirectory(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in ListDirectoryInput,
) (*mcp.CallToolResult, ListDirectoryOutput, error) {
	if strings.TrimSpace(in.Path) == "" {
		return nil, ListDirectoryOutput{}, fmt.Errorf("path cannot be empty")
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
			"%q is a file, not a directory",
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
