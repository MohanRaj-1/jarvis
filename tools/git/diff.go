package git

import (
	"context"
	"fmt"

	internalgit "jarvis/internal/git"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// DiffInput contains the repository path to inspect.
type DiffInput struct {
	Path string `json:"path" jsonschema:"Absolute or relative path to a Git repository"`
}

// DiffOutput contains the current working tree diff.
type DiffOutput struct {
	Diff string `json:"diff"`
}

// Diff returns the current working tree diff for a Git repository.
func Diff(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in DiffInput,
) (*mcp.CallToolResult, DiffOutput, error) {
	diff, err := internalgit.Diff(in.Path)
	if err != nil {
		return nil, DiffOutput{}, fmt.Errorf("get Git working tree diff: %w", err)
	}

	return nil, DiffOutput{Diff: diff}, nil
}
