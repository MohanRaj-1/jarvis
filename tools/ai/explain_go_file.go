package ai

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	internalai "jarvis/internal/ai"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const defaultGeminiModel = "gemini-3.5-flash"

type ExplainGoFileInput struct {
	Path  string `json:"path" jsonschema:"Absolute or relative path to a Go source file"`
	Model string `json:"model,omitempty" jsonschema:"Optional Gemini model name; defaults to gemini-3.5-flash"`
}

type ExplainGoFileOutput struct {
	Explanation string `json:"explanation"`
}

func ExplainGoFile(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in ExplainGoFileInput,
) (*mcp.CallToolResult, ExplainGoFileOutput, error) {
	if strings.TrimSpace(in.Path) == "" {
		return nil, ExplainGoFileOutput{}, fmt.Errorf("path cannot be empty")
	}

	path := filepath.Clean(in.Path)
	if !strings.EqualFold(filepath.Ext(path), ".go") {
		return nil, ExplainGoFileOutput{}, fmt.Errorf("%q is not a Go source file", in.Path)
	}

	model := strings.TrimSpace(in.Model)
	if model == "" {
		model = defaultGeminiModel
	}

	client, err := internalai.NewClient(ctx, internalai.Config{
		Provider: internalai.ProviderGemini,
		APIKey:   os.Getenv("GEMINI_API_KEY"),
		Model:    model,
	})
	if err != nil {
		return nil, ExplainGoFileOutput{}, fmt.Errorf("create AI client: %w", err)
	}

	explanation, err := internalai.ExplainFile(ctx, client, path)
	if err != nil {
		return nil, ExplainGoFileOutput{}, fmt.Errorf("explain Go file: %w", err)
	}

	return nil, ExplainGoFileOutput{Explanation: explanation}, nil
}
