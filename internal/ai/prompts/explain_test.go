package prompts

import (
	"strings"
	"testing"

	"jarvis/internal/analyzer"
)

func TestExplainGoFilePrompt(t *testing.T) {
	prompt := ExplainGoFilePrompt(&analyzer.Analysis{}, "package example")

	for _, want := range []string{"# Purpose", "# Possible Improvements", "package example"} {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt does not contain %q", want)
		}
	}
}
