package git

import (
	"context"
	"fmt"

	aigit "jarvis/internal/ai/git"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ReleaseNotesInput contains a Git repository path and revision range.
type ReleaseNotesInput struct {
	Path string `json:"path" jsonschema:"Absolute or relative path to a Git repository"`
	From string `json:"from" jsonschema:"Starting Git revision to exclude from the range"`
	To   string `json:"to" jsonschema:"Ending Git revision to include in the range, such as HEAD"`
}

// ReleaseNotesOutput contains structured release notes for a commit range.
type ReleaseNotesOutput struct {
	Summary         string   `json:"summary"`
	Features        []string `json:"features"`
	Fixes           []string `json:"fixes"`
	Changes         []string `json:"changes"`
	BreakingChanges []string `json:"breaking_changes"`
}

// NewGenerateReleaseNotes returns an MCP tool backed by service.
func NewGenerateReleaseNotes(service aigit.ReleaseNotesService) func(context.Context, *mcp.CallToolRequest, ReleaseNotesInput) (*mcp.CallToolResult, ReleaseNotesOutput, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, in ReleaseNotesInput) (*mcp.CallToolResult, ReleaseNotesOutput, error) {
		notes, err := service.GenerateReleaseNotes(ctx, in.Path, in.From, in.To)
		if err != nil {
			return nil, ReleaseNotesOutput{}, fmt.Errorf("generate Git release notes: %w", err)
		}
		return nil, ReleaseNotesOutput{
			Summary: notes.Summary, Features: notes.Features, Fixes: notes.Fixes, Changes: notes.Changes, BreakingChanges: notes.BreakingChanges,
		}, nil
	}
}
