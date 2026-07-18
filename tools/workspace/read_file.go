package workspace

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const maxReadFileSize int64 = 1 * 1024 * 1024

// ReadFileInput contains the path to read.
type ReadFileInput struct {
	Path string `json:"path" jsonschema:"Absolute or relative path to the file"`
}

// ReadFileOutput contains file content and its byte size.
type ReadFileOutput struct {
	Content string `json:"content"`
	Size    int64  `json:"size"`
}

// ReadFile reads a file up to the configured size limit.
func ReadFile(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in ReadFileInput,
) (*mcp.CallToolResult, ReadFileOutput, error) {

	if strings.TrimSpace(in.Path) == "" {
		return nil, ReadFileOutput{}, fmt.Errorf("path is required; provide a path to a readable file")
	}
	cleanPath := filepath.Clean(in.Path)

	info, err := os.Stat(cleanPath)
	if err != nil {
		return nil, ReadFileOutput{}, fmt.Errorf("cannot access %q: %w", in.Path, err)
	}

	if info.IsDir() {
		return nil, ReadFileOutput{}, fmt.Errorf("%q is a directory, not a file; provide a file path", in.Path)
	}

	if info.Size() > maxReadFileSize {
		return nil, ReadFileOutput{}, fmt.Errorf(
			"file %q is too large: %d bytes exceeds the 1 MB limit",
			in.Path,
			info.Size(),
		)
	}

	content, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, ReadFileOutput{}, fmt.Errorf("cannot read %q: %w", in.Path, err)
	}

	size := int64(len(content))
	if size > maxReadFileSize {
		return nil, ReadFileOutput{}, fmt.Errorf(
			"file %q is too large: %d bytes exceeds the 1 MB limit",
			in.Path,
			size,
		)
	}

	return nil, ReadFileOutput{
		Content: string(content),
		Size:    size,
	}, nil
}
