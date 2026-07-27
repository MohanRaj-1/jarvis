package git

import (
	"context"
	"fmt"

	aigit "jarvis/internal/ai/git"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GenerateCommitMessageInput contains the repository path to inspect.
type GenerateCommitMessageInput struct {
	Path string `json:"path" jsonschema:"Absolute or relative path to a Git repository"`
}

// GenerateCommitMessageOutput contains the generated Conventional Commit message.
type GenerateCommitMessageOutput struct {
	Type    string `json:"type"`
	Scope   string `json:"scope,omitempty"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// NewGenerateCommitMessage returns an MCP tool backed by service.
func NewGenerateCommitMessage(service aigit.CommitMessageService) func(context.Context, *mcp.CallToolRequest, GenerateCommitMessageInput) (*mcp.CallToolResult, GenerateCommitMessageOutput, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, in GenerateCommitMessageInput) (*mcp.CallToolResult, GenerateCommitMessageOutput, error) {
		message, err := service.GenerateCommitMessage(ctx, in.Path)
		if err != nil {
			return nil, GenerateCommitMessageOutput{}, fmt.Errorf("generate Git commit message: %w", err)
		}

		return nil, GenerateCommitMessageOutput{
			Type:    message.Type,
			Scope:   message.Scope,
			Subject: message.Subject,
			Message: message.String(),
		}, nil
	}
}
