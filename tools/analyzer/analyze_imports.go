package analyzer

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	internalanalyzer "jarvis/internal/analyzer"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type AnalyzeImportsInput struct {
	Path string `json:"path" jsonschema:"Absolute or relative path to a Go source file"`
}

type AnalyzeImportsOutput struct {
	Imports []string `json:"imports"`
}

func AnalyzeImports(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in AnalyzeImportsInput,
) (*mcp.CallToolResult, AnalyzeImportsOutput, error) {
	if strings.TrimSpace(in.Path) == "" {
		return nil, AnalyzeImportsOutput{}, fmt.Errorf("path cannot be empty")
	}

	imports, err := internalanalyzer.ExtractImports(filepath.Clean(in.Path))
	if err != nil {
		return nil, AnalyzeImportsOutput{}, fmt.Errorf(
			"failed to analyze imports: %w",
			err,
		)
	}

	return nil, AnalyzeImportsOutput{
		Imports: imports,
	}, nil
}
