package analyzer

import (
	"context"
	"fmt"

	internalanalyzer "jarvis/internal/analyzer"
	"jarvis/internal/gofile"
	internalworkspace "jarvis/internal/workspace"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// AnalyzeGoFileInput contains the Go source file to analyze.
type AnalyzeGoFileInput struct {
	Path string `json:"path" jsonschema:"Absolute or relative path to a Go source file"`
}

// AnalyzeGoFileOutput contains structural information extracted from a Go file.
type AnalyzeGoFileOutput struct {
	Package    string                       `json:"package"`
	Imports    []string                     `json:"imports"`
	Functions  []string                     `json:"functions"`
	Structs    []string                     `json:"structs"`
	Methods    []internalanalyzer.Method    `json:"methods"`
	Interfaces []internalanalyzer.Interface `json:"interfaces"`
	Todos      []internalanalyzer.Todo      `json:"todos"`
}

// AnalyzeGoFileTool analyzes Go source files in its workspace.
type AnalyzeGoFileTool struct {
	workspace internalworkspace.Workspace
}

// NewAnalyzeGoFileTool creates a Go-file analyzer limited to w.
func NewAnalyzeGoFileTool(w internalworkspace.Workspace) AnalyzeGoFileTool {
	return AnalyzeGoFileTool{workspace: w}
}

// AnalyzeGoFile returns structural information for a Go source file.
func (t AnalyzeGoFileTool) Handle(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in AnalyzeGoFileInput,
) (*mcp.CallToolResult, AnalyzeGoFileOutput, error) {
	path, err := t.workspace.Resolve(in.Path)
	if err != nil {
		return nil, AnalyzeGoFileOutput{}, err
	}

	path, err = gofile.ValidatePath(path)
	if err != nil {
		return nil, AnalyzeGoFileOutput{}, err
	}

	analysis, err := internalanalyzer.Analyze(path)
	if err != nil {
		return nil, AnalyzeGoFileOutput{}, fmt.Errorf(
			"analyze Go source file %q: %w",
			path,
			err,
		)
	}

	return nil, AnalyzeGoFileOutput{
		Package:    analysis.Package,
		Imports:    analysis.Imports,
		Functions:  analysis.Functions,
		Structs:    analysis.Structs,
		Methods:    analysis.Methods,
		Interfaces: analysis.Interfaces,
		Todos:      analysis.Todos,
	}, nil
}
