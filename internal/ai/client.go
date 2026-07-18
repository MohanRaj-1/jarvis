package ai

import (
	"context"
)

// Client generates content from prompts.
type Client interface {
	Generate(ctx context.Context, prompt string) (string, error)
}
