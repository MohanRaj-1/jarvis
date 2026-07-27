package git_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"jarvis/internal/ai/git"
)

type fakeRepository struct {
	diff string
	err  error
}

func (r fakeRepository) Diff(string) (string, error) { return r.diff, r.err }

type fakeAI struct {
	message string
	err     error
	prompt  string
}

func (a *fakeAI) Generate(_ context.Context, prompt string) (string, error) {
	a.prompt = prompt
	return a.message, a.err
}

func TestCommitMessageServiceGenerate(t *testing.T) {
	ai := &fakeAI{message: `{"type":"feat","scope":"api","subject":"add authentication"}`}
	service := git.CommitMessageService{
		Git: fakeRepository{diff: "diff --git a/api.go b/api.go\n+// authentication"},
		AI:  ai,
	}

	message, err := service.GenerateCommitMessage(context.Background(), "/repo")
	if err != nil {
		t.Fatalf("GenerateCommitMessage() error = %v", err)
	}
	if got := message.String(); got != "feat(api): add authentication" {
		t.Errorf("GenerateCommitMessage().String() = %q", got)
	}
	if !strings.Contains(ai.prompt, "diff --git a/api.go b/api.go") {
		t.Errorf("prompt does not include the Git diff: %q", ai.prompt)
	}
}

func TestCommitMessageServiceGenerateWithoutScope(t *testing.T) {
	service := git.CommitMessageService{
		Git: fakeRepository{diff: "diff"},
		AI:  &fakeAI{message: `{"type":"docs","scope":"","subject":"update README"}`},
	}

	message, err := service.GenerateCommitMessage(context.Background(), "/repo")
	if err != nil {
		t.Fatalf("GenerateCommitMessage() error = %v", err)
	}
	if got := message.String(); got != "docs: update README" {
		t.Errorf("GenerateCommitMessage().String() = %q", got)
	}
}

func TestCommitMessageServiceGenerateErrors(t *testing.T) {
	tests := []struct {
		name    string
		service git.CommitMessageService
		want    string
	}{
		{name: "missing Git repository", service: git.CommitMessageService{AI: &fakeAI{}}, want: "Git repository is required"},
		{name: "missing AI client", service: git.CommitMessageService{Git: fakeRepository{}}, want: "AI client is required"},
		{name: "empty diff", service: git.CommitMessageService{Git: fakeRepository{}, AI: &fakeAI{}}, want: "no uncommitted changes"},
		{name: "Git failure", service: git.CommitMessageService{Git: fakeRepository{err: errors.New("failed")}, AI: &fakeAI{}}, want: "get Git diff"},
		{name: "AI failure", service: git.CommitMessageService{Git: fakeRepository{diff: "diff"}, AI: &fakeAI{err: errors.New("failed")}}, want: "generate commit message"},
		{name: "invalid AI response", service: git.CommitMessageService{Git: fakeRepository{diff: "diff"}, AI: &fakeAI{message: "feat: add feature"}}, want: "parse generated commit message"},
		{name: "multi-line subject", service: git.CommitMessageService{Git: fakeRepository{diff: "diff"}, AI: &fakeAI{message: `{"type":"feat","scope":"ai","subject":"add feature\nwith details"}`}}, want: "single-line"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.service.GenerateCommitMessage(context.Background(), "/repo")
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("Generate() error = %v, want %q", err, tt.want)
			}
		})
	}
}
