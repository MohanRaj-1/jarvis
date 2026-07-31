package git_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"jarvis/internal/ai/git"
	internalgit "jarvis/internal/git"
)

func TestExplainCommitServiceExplainCommit(t *testing.T) {
	ai := &fakeAI{message: "This adds the branch-listing capability and exposes it through MCP."}
	service := git.ExplainCommitService{
		Git: fakeRepository{commit: &internalgit.CommitDetails{
			Commit: internalgit.Commit{Message: "feat(git): add repository branches tool", Author: "Mohan Raj"},
			Files: []internalgit.ChangedFile{
				{Path: "internal/git/branches.go", Status: "Added"},
				{Path: "tools/git/branches.go", Status: "Added"},
			},
		}},
		AI: ai,
	}

	summary, err := service.ExplainCommit(context.Background(), "/repo", "91b6cd2")
	if err != nil {
		t.Fatalf("ExplainCommit() error = %v", err)
	}
	if summary != "This adds the branch-listing capability and exposes it through MCP." {
		t.Errorf("ExplainCommit() = %q", summary)
	}
	for _, want := range []string{"feat(git): add repository branches tool", "Mohan Raj", "internal/git/branches.go"} {
		if !strings.Contains(ai.prompt, want) {
			t.Errorf("prompt does not contain %q: %q", want, ai.prompt)
		}
	}
}

func TestExplainCommitServiceErrors(t *testing.T) {
	tests := []struct {
		name    string
		service git.ExplainCommitService
		want    string
	}{
		{name: "missing Git repository", service: git.ExplainCommitService{AI: &fakeAI{}}, want: "Git repository is required"},
		{name: "missing AI client", service: git.ExplainCommitService{Git: fakeRepository{}}, want: "AI client is required"},
		{name: "unknown hash", service: git.ExplainCommitService{Git: fakeRepository{showErr: errors.New("commit not found")}, AI: &fakeAI{}}, want: "show Git commit"},
		{name: "Git error", service: git.ExplainCommitService{Git: fakeRepository{showErr: errors.New("repository unavailable")}, AI: &fakeAI{}}, want: "show Git commit"},
		{name: "AI error", service: git.ExplainCommitService{Git: fakeRepository{commit: &internalgit.CommitDetails{}}, AI: &fakeAI{err: errors.New("failed")}}, want: "generate commit explanation"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.service.ExplainCommit(context.Background(), "/repo", "91b6cd2")
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("ExplainCommit() error = %v, want %q", err, tt.want)
			}
		})
	}
}
