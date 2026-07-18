package workspace

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// FindFilesInput contains the search directory and filepath glob.
type FindFilesInput struct {
	Root    string `json:"root" jsonschema:"Absolute or relative directory to search recursively"`
	Pattern string `json:"pattern" jsonschema:"Filepath glob pattern, for example *.go or tools/*.go"`
}

// FindFilesOutput contains paths that matched the search pattern.
type FindFilesOutput struct {
	Paths []string `json:"paths"`
}

// FindFiles recursively finds files matching a filepath glob.
func FindFiles(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in FindFilesInput,
) (*mcp.CallToolResult, FindFilesOutput, error) {
	if strings.TrimSpace(in.Root) == "" {
		return nil, FindFilesOutput{}, fmt.Errorf("root is required; provide a directory to search")
	}
	if strings.TrimSpace(in.Pattern) == "" {
		return nil, FindFilesOutput{}, fmt.Errorf("pattern is required; provide a filepath glob such as *.go")
	}

	root := filepath.Clean(in.Root)
	pattern := filepath.FromSlash(in.Pattern)
	info, err := os.Stat(root)
	if err != nil {
		return nil, FindFilesOutput{}, fmt.Errorf("cannot access root %q: %w", root, err)
	}
	if !info.IsDir() {
		return nil, FindFilesOutput{}, fmt.Errorf("%q is a file, not a directory; provide a directory to search", root)
	}

	if _, err := filepath.Match(pattern, ""); err != nil {
		return nil, FindFilesOutput{}, fmt.Errorf("invalid pattern %q: %w", in.Pattern, err)
	}

	paths := make([]string, 0)
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}

		relativePath, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("cannot make %q relative to root %q: %w", path, root, err)
		}

		nameMatches, err := filepath.Match(pattern, entry.Name())
		if err != nil {
			return err
		}
		pathMatches, err := filepath.Match(pattern, relativePath)
		if err != nil {
			return err
		}
		if nameMatches || pathMatches {
			paths = append(paths, relativePath)
		}

		return nil
	})
	if err != nil {
		return nil, FindFilesOutput{}, fmt.Errorf("cannot search root %q: %w", root, err)
	}

	return nil, FindFilesOutput{Paths: paths}, nil
}
