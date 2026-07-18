package ai

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

// GeminiClient generates content with the Google Gemini API.
type GeminiClient struct {
	client *genai.Client
	model  string
}

// NewGeminiClient creates a Gemini client from the supplied configuration.
func NewGeminiClient(ctx context.Context, cfg Config) (*GeminiClient, error) {
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, errors.New("Gemini API key is required; set the GEMINI_API_KEY environment variable")
	}
	if strings.TrimSpace(cfg.Model) == "" {
		return nil, errors.New("Gemini model is required; provide a Gemini model name")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: cfg.APIKey})
	if err != nil {
		return nil, fmt.Errorf("initialize Gemini client: %w", err)
	}

	return &GeminiClient{client: client, model: cfg.Model}, nil
}

// Generate sends a prompt to Gemini and returns the generated text.
func (c *GeminiClient) Generate(ctx context.Context, prompt string) (string, error) {
	response, err := c.client.Models.GenerateContent(ctx, c.model, genai.Text(prompt), nil)
	if err != nil {
		return "", fmt.Errorf("generate content with Gemini model %q: %w", c.model, err)
	}

	return response.Text(), nil
}
