package ai

import (
	"context"
)

// Client represents an AI provider.
type Client interface {
	Generate(ctx context.Context, prompt string) (string, error)
}
