package ai

import (
	"context"
	"errors"
	"fmt"
	"os"

	"jarvis/internal/analyzer"
)

// ExplainFile analyzes a Go source file and generates an AI explanation for it.
func ExplainFile(ctx context.Context, client Client, path string) (string, error) {
	if client == nil {
		return "", errors.New("AI client is required to explain a file")
	}

	source, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read Go source file %q: %w", path, err)
	}

	analysis, err := analyzer.Analyze(path)
	if err != nil {
		return "", fmt.Errorf("analyze Go source file %q: %w", path, err)
	}

	prompt := BuildExplainPrompt(analysis, string(source))
	explanation, err := client.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("generate explanation for Go source file %q: %w", path, err)
	}

	return explanation, nil
}
