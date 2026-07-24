package git

import (
	"context"
	"fmt"
	"time"

	internalgit "jarvis/internal/git"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ShowCommitInput contains the repository path and commit hash to inspect.
type ShowCommitInput struct {
	Path string `json:"path" jsonschema:"Absolute or relative path to a Git repository"`
	Hash string `json:"hash" jsonschema:"Commit hash, including an unambiguous abbreviated hash"`
}

// ChangedFileOutput describes a file changed by a commit.
type ChangedFileOutput struct {
	Path   string `json:"path"`
	Status string `json:"status"`
}

// ShowCommitOutput contains detailed information about a commit.
type ShowCommitOutput struct {
	Hash    string              `json:"hash"`
	Author  string              `json:"author"`
	Message string              `json:"message"`
	Date    time.Time           `json:"date"`
	Parents []string            `json:"parents"`
	Files   []ChangedFileOutput `json:"files"`
}

// ShowCommit returns detailed information about a Git commit.
func ShowCommit(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in ShowCommitInput,
) (*mcp.CallToolResult, ShowCommitOutput, error) {
	details, err := internalgit.Show(in.Path, in.Hash)
	if err != nil {
		return nil, ShowCommitOutput{}, fmt.Errorf("show Git commit: %w", err)
	}

	output := ShowCommitOutput{
		Hash:    details.Hash,
		Author:  details.Author,
		Message: details.Message,
		Date:    details.Date,
		Parents: details.Parents,
		Files:   make([]ChangedFileOutput, len(details.Files)),
	}
	for i, file := range details.Files {
		output.Files[i] = ChangedFileOutput{Path: file.Path, Status: file.Status}
	}
	return nil, output, nil
}
