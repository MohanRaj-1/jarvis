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

func TestReleaseNotesPrompt(t *testing.T) {
	prompt := ReleaseNotesPrompt("v0.5.2", "HEAD", []internalgit.ReleaseCommit{{Hash: "abcdef0", Author: "Mohan Raj", Date: time.Date(2026, time.August, 6, 0, 0, 0, 0, time.UTC), Message: "feat(git): add release notes", ChangedFiles: []internalgit.ChangedFile{{Path: "internal/ai/git/release_notes.go", Status: "Added"}}}})

	for _, want := range []string{"Return valid JSON only.", `"summary"`, `"features"`, `"fixes"`, `"changes"`, `"breaking_changes"`, "v0.5.2", "HEAD", "abcdef0", "Mohan Raj", "2026-08-06T00:00:00Z", "feat(git): add release notes", "internal/ai/git/release_notes.go", "Added"} {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt does not contain %q", want)
		}
	}
}
