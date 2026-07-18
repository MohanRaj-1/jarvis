package ai

import (
	"context"
	"fmt"
)

// NewClient creates an AI client for the configured provider.
func NewClient(ctx context.Context, config Config) (Client, error) {
	switch config.Provider {
	case ProviderGemini:
		return NewGeminiClient(ctx, config)
	default:
		return nil, fmt.Errorf("unsupported AI provider %q; configure a supported provider such as %q", config.Provider, ProviderGemini)
	}
}
