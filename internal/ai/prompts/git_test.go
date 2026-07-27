package prompts

import (
	"strings"
	"testing"
)

func TestCommitMessagePrompt(t *testing.T) {
	diff := "diff --git a/a.go b/a.go"
	prompt := CommitMessagePrompt(diff)

	for _, want := range []string{"Conventional Commit", "Return valid JSON only.", `"type"`, `"scope"`, `"subject"`, diff} {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt does not contain %q", want)
		}
	}
}
