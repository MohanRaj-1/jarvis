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

// BranchesTool returns local branches from its repository.
type BranchesTool struct {
	repository internalgit.DefaultRepository
}

// NewBranchesTool creates a branches tool using repository.
func NewBranchesTool(repository internalgit.DefaultRepository) BranchesTool {
	return BranchesTool{repository: repository}
}

// RepositoryBranches returns the current branch and all local branches for a Git repository.
func (t BranchesTool) Handle(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in BranchesInput,
) (*mcp.CallToolResult, BranchesOutput, error) {
	branches, err := t.repository.RepositoryBranches(in.Path)
	if err != nil {
		return nil, BranchesOutput{}, fmt.Errorf("get Git repository branches: %w", err)
	}

	return nil, BranchesOutput{
		Current:  branches.Current,
		Branches: branches.Local,
	}, nil
}
