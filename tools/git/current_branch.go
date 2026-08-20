// Package git contains MCP tools related to Git repositories.
package git

import (
	"context"
	"fmt"

	internalgit "jarvis/internal/git"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// CurrentBranchInput contains the repository path to inspect.
type CurrentBranchInput struct {
	RepoPath string `json:"repo_path" jsonschema:"Absolute or relative path to a Git repository"`
}

// CurrentBranchOutput contains the currently checked-out branch name.
type CurrentBranchOutput struct {
	Branch string `json:"branch"`
}

// CurrentBranchTool returns the current branch of its repository.
type CurrentBranchTool struct {
	repository internalgit.DefaultRepository
}

// NewCurrentBranchTool creates a current-branch tool using repository.
func NewCurrentBranchTool(repository internalgit.DefaultRepository) CurrentBranchTool {
	return CurrentBranchTool{repository: repository}
}

// CurrentBranch returns the current branch for a Git repository.
func (t CurrentBranchTool) Handle(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in CurrentBranchInput,
) (*mcp.CallToolResult, CurrentBranchOutput, error) {
	branch, err := t.repository.CurrentBranch(in.RepoPath)
	if err != nil {
		return nil, CurrentBranchOutput{}, fmt.Errorf("get current Git branch: %w", err)
	}

	return nil, CurrentBranchOutput{Branch: branch}, nil
}
