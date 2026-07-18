package analyzer

import (
	"context"
	"fmt"

	internalanalyzer "jarvis/internal/analyzer"
	"jarvis/internal/gofile"

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

// AnalyzeGoFile returns structural information for a Go source file.
func AnalyzeGoFile(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in AnalyzeGoFileInput,
) (*mcp.CallToolResult, AnalyzeGoFileOutput, error) {
	path, err := gofile.ValidatePath(in.Path)
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
