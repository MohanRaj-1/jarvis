package analyzer

import (
	"context"
	"fmt"

	internalanalyzer "jarvis/internal/analyzer"
	"jarvis/internal/gofile"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// AnalyzeImportsInput contains the Go source file to inspect.
type AnalyzeImportsInput struct {
	Path string `json:"path" jsonschema:"Absolute or relative path to a Go source file"`
}

// AnalyzeImportsOutput contains the import paths declared by a Go file.
type AnalyzeImportsOutput struct {
	Imports []string `json:"imports"`
}

// AnalyzeImports returns import paths from a Go source file.
func AnalyzeImports(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in AnalyzeImportsInput,
) (*mcp.CallToolResult, AnalyzeImportsOutput, error) {
	path, err := gofile.ValidatePath(in.Path)
	if err != nil {
		return nil, AnalyzeImportsOutput{}, err
	}

	imports, err := internalanalyzer.ExtractImports(path)
	if err != nil {
		return nil, AnalyzeImportsOutput{}, fmt.Errorf(
			"analyze imports in Go source file %q: %w",
			path,
			err,
		)
	}

	return nil, AnalyzeImportsOutput{
		Imports: imports,
	}, nil
}
