package analyzer

import (
	"context"
	"fmt"

	internalanalyzer "jarvis/internal/analyzer"
	"jarvis/internal/gofile"
	internalworkspace "jarvis/internal/workspace"

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

// AnalyzeImportsTool extracts imports from Go source files in its workspace.
type AnalyzeImportsTool struct {
	workspace internalworkspace.Workspace
}

// NewAnalyzeImportsTool creates an import analyzer limited to w.
func NewAnalyzeImportsTool(w internalworkspace.Workspace) AnalyzeImportsTool {
	return AnalyzeImportsTool{workspace: w}
}

// AnalyzeImports returns import paths from a Go source file.
func (t AnalyzeImportsTool) Handle(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in AnalyzeImportsInput,
) (*mcp.CallToolResult, AnalyzeImportsOutput, error) {
	path, err := t.workspace.Resolve(in.Path)
	if err != nil {
		return nil, AnalyzeImportsOutput{}, err
	}

	path, err = gofile.ValidatePath(path)
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
