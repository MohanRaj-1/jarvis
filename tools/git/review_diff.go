package git

import (
	"context"
	"fmt"

	aigit "jarvis/internal/ai/git"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ReviewDiffInput contains the repository path to review.
type ReviewDiffInput struct {
	Path string `json:"path" jsonschema:"Absolute or relative path to a Git repository"`
}

// ReviewDiffOutput contains the structured Git diff review.
type ReviewDiffOutput struct {
	Summary      string   `json:"summary"`
	Strengths    []string `json:"strengths"`
	Issues       []string `json:"issues"`
	Suggestions  []string `json:"suggestions"`
	OverallScore int      `json:"overall_score"`
}

// NewReviewDiff returns an MCP tool backed by service.
func NewReviewDiff(service aigit.ReviewDiffService) func(context.Context, *mcp.CallToolRequest, ReviewDiffInput) (*mcp.CallToolResult, ReviewDiffOutput, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, in ReviewDiffInput) (*mcp.CallToolResult, ReviewDiffOutput, error) {
		report, err := service.ReviewDiff(ctx, in.Path)
		if err != nil {
			return nil, ReviewDiffOutput{}, fmt.Errorf("review Git diff: %w", err)
		}

		return nil, ReviewDiffOutput{
			Summary:      report.Summary,
			Strengths:    report.Strengths,
			Issues:       report.Issues,
			Suggestions:  report.Suggestions,
			OverallScore: report.OverallScore,
		}, nil
	}
}
