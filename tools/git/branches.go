package git

import (
	"context"
	"fmt"

	internalgit "jarvis/internal/git"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// BranchesInput contains the repository path to inspect.
type BranchesInput struct {
	Path string `json:"path" jsonschema:"Absolute or relative path to a Git repository"`
}

// BranchesOutput contains the current branch and all local branches.
type BranchesOutput struct {
	Current  string   `json:"current"`
	Branches []string `json:"branches"`
}

// RepositoryBranches returns the current branch and all local branches for a Git repository.
func RepositoryBranches(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in BranchesInput,
) (*mcp.CallToolResult, BranchesOutput, error) {
	branches, err := internalgit.RepositoryBranches(in.Path)
	if err != nil {
		return nil, BranchesOutput{}, fmt.Errorf("get Git repository branches: %w", err)
	}

	return nil, BranchesOutput{
		Current:  branches.Current,
		Branches: branches.Local,
	}, nil
}
