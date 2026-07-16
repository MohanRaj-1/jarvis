package ai

import (
	"context"
	"fmt"
)

func NewClient(ctx context.Context, config Config) (Client, error) {
	switch config.Provider {
	case ProviderGemini:
		return NewGeminiClient(ctx, config)
	default:
		return nil, fmt.Errorf("unsupported provider %q", config.Provider)
	}
}
