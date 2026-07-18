package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

	internalai "jarvis/internal/ai"
	"jarvis/internal/gofile"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const defaultGeminiModel = "gemini-3.5-flash"

// ExplainGoFileInput contains the Go source file and optional Gemini model to use.
type ExplainGoFileInput struct {
	Path  string `json:"path" jsonschema:"Absolute or relative path to a Go source file"`
	Model string `json:"model,omitempty" jsonschema:"Optional Gemini model name; defaults to gemini-3.5-flash"`
}

// ExplainGoFileOutput contains the generated source-code explanation.
type ExplainGoFileOutput struct {
	Explanation string `json:"explanation"`
}

// ExplainGoFile explains a Go source file with Gemini.
func ExplainGoFile(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in ExplainGoFileInput,
) (*mcp.CallToolResult, ExplainGoFileOutput, error) {
	path, err := gofile.ValidatePath(in.Path)
	if err != nil {
		return nil, ExplainGoFileOutput{}, err
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
		return nil, ExplainGoFileOutput{}, fmt.Errorf("configure Gemini for Go file explanation: %w", err)
	}

	explanation, err := internalai.ExplainFile(ctx, client, path)
	if err != nil {
		return nil, ExplainGoFileOutput{}, fmt.Errorf("explain Go source file %q: %w", path, err)
	}

	return nil, ExplainGoFileOutput{Explanation: explanation}, nil
}
