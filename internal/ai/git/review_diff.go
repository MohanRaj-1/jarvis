package git

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"jarvis/internal/ai"
	"jarvis/internal/ai/prompts"
	internalgit "jarvis/internal/git"
)

// ReviewReport contains a structured review of a Git working tree diff.
type ReviewReport struct {
	Summary      string   `json:"summary"`
	Strengths    []string `json:"strengths"`
	Issues       []string `json:"issues"`
	Suggestions  []string `json:"suggestions"`
	OverallScore int      `json:"overall_score"`
}

// ReviewDiffService generates structured code reviews from Git diffs.
type ReviewDiffService struct {
	Git internalgit.Repository
	AI  ai.Client
}

type reviewResponse struct {
	Summary      string   `json:"summary"`
	Strengths    []string `json:"strengths"`
	Issues       []string `json:"issues"`
	Suggestions  []string `json:"suggestions"`
	OverallScore int      `json:"overall_score"`
}

// ReviewDiff reviews uncommitted changes in repoPath.
func (s ReviewDiffService) ReviewDiff(ctx context.Context, repoPath string) (*ReviewReport, error) {
	if s.Git == nil {
		return nil, errors.New("Git repository is required to review a diff")
	}
	if s.AI == nil {
		return nil, errors.New("AI client is required to review a diff")
	}

	diff, err := s.Git.Diff(repoPath)
	if err != nil {
		return nil, fmt.Errorf("get Git diff: %w", err)
	}
	if strings.TrimSpace(diff) == "" {
		return nil, errors.New("cannot review a diff: repository has no uncommitted changes")
	}

	response, err := s.AI.Generate(ctx, prompts.ReviewDiffPrompt(diff))
	if err != nil {
		return nil, fmt.Errorf("generate diff review: %w", err)
	}

	report, err := parseReviewReport(response)
	if err != nil {
		return nil, fmt.Errorf("parse generated diff review: %w", err)
	}

	return report, nil
}

func parseReviewReport(response string) (*ReviewReport, error) {
	var review reviewResponse
	if err := json.Unmarshal([]byte(response), &review); err != nil {
		return nil, fmt.Errorf("expected JSON review report: %w", err)
	}

	review.Summary = strings.TrimSpace(review.Summary)
	if review.Summary == "" {
		return nil, errors.New("summary is required")
	}
	if review.OverallScore < 1 || review.OverallScore > 10 {
		return nil, errors.New("overall_score must be between 1 and 10")
	}

	return &ReviewReport{
		Summary:      review.Summary,
		Strengths:    review.Strengths,
		Issues:       review.Issues,
		Suggestions:  review.Suggestions,
		OverallScore: review.OverallScore,
	}, nil
}
