package workspace

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// FileInfoInput contains the path to inspect.
type FileInfoInput struct {
	Path string `json:"path" jsonschema:"Absolute or relative path to a file or directory"`
}

// FileInfoOutput contains metadata for a filesystem path.
type FileInfoOutput struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	IsDir    bool   `json:"is_dir"`
	Modified string `json:"modified"`
}

// FileInfo returns metadata for a file or directory.
func FileInfo(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in FileInfoInput,
) (*mcp.CallToolResult, FileInfoOutput, error) {
	if strings.TrimSpace(in.Path) == "" {
		return nil, FileInfoOutput{}, fmt.Errorf("path is required; provide a file or directory path")
	}

	cleanPath := filepath.Clean(in.Path)
	info, err := os.Stat(cleanPath)
	if err != nil {
		return nil, FileInfoOutput{}, fmt.Errorf("cannot access %q: %w", cleanPath, err)
	}

	return nil, FileInfoOutput{
		Name:     info.Name(),
		Size:     info.Size(),
		IsDir:    info.IsDir(),
		Modified: info.ModTime().Format(time.RFC3339),
	}, nil
}
