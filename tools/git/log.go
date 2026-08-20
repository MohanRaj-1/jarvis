package git

import (
	"context"
	"fmt"
	"time"

	internalgit "jarvis/internal/git"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// LogInput contains the repository path and optional number of commits to return.
type LogInput struct {
	Path  string `json:"path" jsonschema:"Absolute or relative path to a Git repository"`
	Limit int    `json:"limit,omitempty" jsonschema:"Maximum number of commits to return; defaults to 10"`
}

// CommitOutput represents a commit returned by the git_log tool.
type CommitOutput struct {
	Hash    string    `json:"hash"`
	Author  string    `json:"author"`
	Email   string    `json:"email"`
	Date    time.Time `json:"date"`
	Message string    `json:"message"`
}

// LogOutput contains the requested commits.
type LogOutput struct {
	Commits []CommitOutput `json:"commits"`
}

// Log returns the most recent commits for a Git repository.
func Log(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in LogInput,
) (*mcp.CallToolResult, LogOutput, error) {
	commits, err := internalgit.Log(in.Path, in.Limit)
	if err != nil {
		return nil, LogOutput{}, fmt.Errorf("get Git repository log: %w", err)
	}

	output := LogOutput{Commits: make([]CommitOutput, len(commits))}
	for i, commit := range commits {
		output.Commits[i] = CommitOutput{
			Hash:    commit.Hash,
			Author:  commit.Author,
			Email:   commit.Email,
			Date:    commit.Date,
			Message: commit.Message,
		}
	}
	return nil, output, nil
}
