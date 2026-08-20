package workspace

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	internalworkspace "jarvis/internal/workspace"

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

// FileInfoTool returns metadata for paths in its workspace.
type FileInfoTool struct {
	workspace internalworkspace.Workspace
}

// NewFileInfoTool creates a file-info tool limited to w.
func NewFileInfoTool(w internalworkspace.Workspace) FileInfoTool {
	return FileInfoTool{workspace: w}
}

// FileInfo returns metadata for a file or directory.
func (t FileInfoTool) Handle(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in FileInfoInput,
) (*mcp.CallToolResult, FileInfoOutput, error) {
	if strings.TrimSpace(in.Path) == "" {
		return nil, FileInfoOutput{}, fmt.Errorf("path is required; provide a file or directory path")
	}

	cleanPath, err := t.workspace.Resolve(in.Path)
	if err != nil {
		return nil, FileInfoOutput{}, err
	}
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
