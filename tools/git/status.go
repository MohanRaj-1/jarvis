package git

import (
	"context"
	"fmt"

	internalgit "jarvis/internal/git"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// StatusInput contains the repository path to inspect.
type StatusInput struct {
	RepoPath string `json:"repo_path" jsonschema:"Absolute or relative path to a Git repository"`
}

// StatusOutput contains the current repository status.
type StatusOutput struct {
	Branch    string   `json:"branch"`
	Modified  []string `json:"modified"`
	Staged    []string `json:"staged"`
	Untracked []string `json:"untracked"`
}

// RepositoryStatus returns the branch and changed files for a Git repository.
func RepositoryStatus(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in StatusInput,
) (*mcp.CallToolResult, StatusOutput, error) {
	status, err := internalgit.RepositoryStatus(in.RepoPath)
	if err != nil {
		return nil, StatusOutput{}, fmt.Errorf("get Git repository status: %w", err)
	}

	return nil, StatusOutput{
		Branch:    status.Branch,
		Modified:  status.Modified,
		Staged:    status.Staged,
		Untracked: status.Untracked,
	}, nil
}
