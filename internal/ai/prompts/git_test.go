package prompts

import (
	"strings"
	"testing"
	"time"

	internalgit "jarvis/internal/git"
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

func TestReviewDiffPrompt(t *testing.T) {
	diff := "diff --git a/a.go b/a.go"
	prompt := ReviewDiffPrompt(diff)

	for _, want := range []string{"senior Go engineer", "Return ONLY valid JSON.", `"summary"`, `"strengths"`, `"issues"`, `"suggestions"`, `"overall_score"`, "Review only the supplied diff.", diff} {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt does not contain %q", want)
		}
	}
}

func TestExplainCommitPrompt(t *testing.T) {
	prompt := ExplainCommitPrompt(&internalgit.CommitDetails{
		Commit: internalgit.Commit{Message: "feat(git): add repository branches tool", Author: "Mohan Raj", Date: time.Now()},
		Files:  []internalgit.ChangedFile{{Path: "internal/git/branches.go", Status: "Added"}},
	})

	for _, want := range []string{"Maximum 150 words.", "Do not speculate", "feat(git): add repository branches tool", "Mohan Raj", "internal/git/branches.go"} {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt does not contain %q", want)
		}
	}
}
