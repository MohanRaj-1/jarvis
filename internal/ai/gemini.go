package ai

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
	model  string
}

func NewGeminiClient(ctx context.Context, cfg Config) (*GeminiClient, error) {
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, errors.New("Gemini API key is required")
	}
	if strings.TrimSpace(cfg.Model) == "" {
		return nil, errors.New("Gemini model is required")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: cfg.APIKey})
	if err != nil {
		return nil, fmt.Errorf("create Gemini client: %w", err)
	}

	return &GeminiClient{client: client, model: cfg.Model}, nil
}

func (c *GeminiClient) Generate(ctx context.Context, prompt string) (string, error) {
	response, err := c.client.Models.GenerateContent(ctx, c.model, genai.Text(prompt), nil)
	if err != nil {
		return "", fmt.Errorf("generate Gemini content: %w", err)
	}

	return response.Text(), nil
}
