package git_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"jarvis/internal/ai/git"
)

func TestReviewDiffServiceReviewDiff(t *testing.T) {
	ai := &fakeAI{message: `{
		"summary":"Adds structured diff reviews.",
		"strengths":["Uses dependency injection"],
		"issues":[],
		"suggestions":["Add retry logic for malformed responses"],
		"overall_score":9
	}`}
	service := git.ReviewDiffService{
		Git: fakeRepository{diff: "diff --git a/review.go b/review.go\n+type ReviewReport struct{}"},
		AI:  ai,
	}

	report, err := service.ReviewDiff(context.Background(), "/repo")
	if err != nil {
		t.Fatalf("ReviewDiff() error = %v", err)
	}
	if report.Summary != "Adds structured diff reviews." {
		t.Errorf("ReviewDiff().Summary = %q", report.Summary)
	}
	if report.OverallScore != 9 {
		t.Errorf("ReviewDiff().OverallScore = %d, want 9", report.OverallScore)
	}
	if len(report.Strengths) != 1 || report.Strengths[0] != "Uses dependency injection" {
		t.Errorf("ReviewDiff().Strengths = %#v", report.Strengths)
	}
	if !strings.Contains(ai.prompt, "diff --git a/review.go b/review.go") {
		t.Errorf("prompt does not include the Git diff: %q", ai.prompt)
	}
}

func TestReviewDiffServiceReviewDiffErrors(t *testing.T) {
	tests := []struct {
		name    string
		service git.ReviewDiffService
		want    string
	}{
		{name: "missing Git repository", service: git.ReviewDiffService{AI: &fakeAI{}}, want: "Git repository is required"},
		{name: "missing AI client", service: git.ReviewDiffService{Git: fakeRepository{}}, want: "AI client is required"},
		{name: "empty diff", service: git.ReviewDiffService{Git: fakeRepository{}, AI: &fakeAI{}}, want: "no uncommitted changes"},
		{name: "Git failure", service: git.ReviewDiffService{Git: fakeRepository{err: errors.New("failed")}, AI: &fakeAI{}}, want: "get Git diff"},
		{name: "AI failure", service: git.ReviewDiffService{Git: fakeRepository{diff: "diff"}, AI: &fakeAI{err: errors.New("failed")}}, want: "generate diff review"},
		{name: "invalid JSON from AI", service: git.ReviewDiffService{Git: fakeRepository{diff: "diff"}, AI: &fakeAI{message: "not JSON"}}, want: "parse generated diff review"},
		{name: "invalid score", service: git.ReviewDiffService{Git: fakeRepository{diff: "diff"}, AI: &fakeAI{message: `{"summary":"summary","overall_score":11}`}}, want: "overall_score must be between 1 and 10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.service.ReviewDiff(context.Background(), "/repo")
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("ReviewDiff() error = %v, want %q", err, tt.want)
			}
		})
	}
}
