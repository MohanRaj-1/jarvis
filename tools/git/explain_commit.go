package git

import (
	"context"
	"fmt"

	aigit "jarvis/internal/ai/git"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ExplainCommitInput contains the repository path and commit hash to explain.
type ExplainCommitInput struct {
	Path string `json:"path" jsonschema:"Absolute or relative path to a Git repository"`
	Hash string `json:"hash" jsonschema:"Commit hash, including an unambiguous abbreviated hash"`
}

// ExplainCommitOutput contains the generated commit explanation.
type ExplainCommitOutput struct {
	Summary string `json:"summary"`
}

// NewExplainCommit returns an MCP tool backed by service.
func NewExplainCommit(service aigit.ExplainCommitService) func(context.Context, *mcp.CallToolRequest, ExplainCommitInput) (*mcp.CallToolResult, ExplainCommitOutput, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, in ExplainCommitInput) (*mcp.CallToolResult, ExplainCommitOutput, error) {
		summary, err := service.ExplainCommit(ctx, in.Path, in.Hash)
		if err != nil {
			return nil, ExplainCommitOutput{}, fmt.Errorf("explain Git commit: %w", err)
		}

		return nil, ExplainCommitOutput{Summary: summary}, nil
	}
}
