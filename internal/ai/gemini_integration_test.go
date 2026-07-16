package ai_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"jarvis/internal/ai"
)

func TestGeminiGenerate(t *testing.T) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("set GEMINI_API_KEY to run the Gemini integration test")
	}

	ctx := context.Background()
	cfg := ai.Config{
		Provider: ai.ProviderGemini,
		APIKey:   apiKey,
		Model:    "gemini-3.5-flash",
	}

	client, err := ai.NewClient(ctx, cfg)
	if err != nil {
		t.Fatalf("create Gemini client: %v", err)
	}

	text, err := client.Generate(ctx, "Introduce yourself in one sentence.")
	if err != nil {
		t.Fatalf("generate text: %v", err)
	}

	fmt.Println(text)
}
