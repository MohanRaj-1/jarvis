package git

import (
	"context"
	"fmt"

	aigit "jarvis/internal/ai/git"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RepositorySummaryInput contains the repository path to summarize.
type RepositorySummaryInput struct {
	Path string `json:"path" jsonschema:"Absolute or relative path to a Git repository"`
}

// RepositorySummaryOutput contains a structured repository summary.
type RepositorySummaryOutput struct {
	Overview        string   `json:"overview"`
	CurrentBranch   string   `json:"current_branch"`
	RecentWork      []string `json:"recent_work"`
	WorkingChanges  []string `json:"working_changes"`
	Recommendations []string `json:"recommendations"`
}

// NewRepositorySummary returns an MCP tool backed by service.
func NewRepositorySummary(service aigit.RepositorySummaryService) func(context.Context, *mcp.CallToolRequest, RepositorySummaryInput) (*mcp.CallToolResult, RepositorySummaryOutput, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, in RepositorySummaryInput) (*mcp.CallToolResult, RepositorySummaryOutput, error) {
		summary, err := service.SummarizeRepository(ctx, in.Path)
		if err != nil {
			return nil, RepositorySummaryOutput{}, fmt.Errorf("summarize Git repository: %w", err)
		}
		return nil, RepositorySummaryOutput{
			Overview: summary.Overview, CurrentBranch: summary.CurrentBranch, RecentWork: summary.RecentWork,
			WorkingChanges: summary.WorkingChanges, Recommendations: summary.Recommendations,
		}, nil
	}
}
