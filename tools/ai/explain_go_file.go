package ai

import (
	"context"
	"fmt"

	internalai "jarvis/internal/ai"
	"jarvis/internal/gofile"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ExplainGoFileInput contains the Go source file to explain.
type ExplainGoFileInput struct {
	Path string `json:"path" jsonschema:"Absolute or relative path to a Go source file"`
}

// ExplainGoFileOutput contains the generated source-code explanation.
type ExplainGoFileOutput struct {
	Explanation string `json:"explanation"`
}

// NewExplainGoFile returns an MCP tool that uses the injected AI client.
func NewExplainGoFile(client internalai.Client) func(context.Context, *mcp.CallToolRequest, ExplainGoFileInput) (*mcp.CallToolResult, ExplainGoFileOutput, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, in ExplainGoFileInput) (*mcp.CallToolResult, ExplainGoFileOutput, error) {
		path, err := gofile.ValidatePath(in.Path)
		if err != nil {
			return nil, ExplainGoFileOutput{}, err
		}

		explanation, err := internalai.ExplainFile(ctx, client, path)
		if err != nil {
			return nil, ExplainGoFileOutput{}, fmt.Errorf("explain Go source file %q: %w", path, err)
		}

		return nil, ExplainGoFileOutput{Explanation: explanation}, nil
	}
}
