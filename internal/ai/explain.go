package ai

import (
	"context"
	"errors"
	"fmt"
	"os"

	"jarvis/internal/analyzer"
)

func ExplainFile(ctx context.Context, client Client, path string) (string, error) {
	if client == nil {
		return "", errors.New("provider not implemented")
	}

	source, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}

	analysis, err := analyzer.Analyze(path)
	if err != nil {
		return "", fmt.Errorf("analyze file: %w", err)
	}

	prompt := BuildExplainPrompt(analysis, string(source))
	explanation, err := client.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("generate explanation: %w", err)
	}

	return explanation, nil
}
