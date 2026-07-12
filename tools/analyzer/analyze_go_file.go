package analyzer

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	internalanalyzer "jarvis/internal/analyzer"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type AnalyzeGoFileInput struct {
	Path string `json:"path" jsonschema:"Absolute or relative path to a Go source file"`
}

type AnalyzeGoFileOutput struct {
	Imports   []string `json:"imports"`
	Functions []string `json:"functions"`
	Structs   []string `json:"structs"`
}

func AnalyzeGoFile(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in AnalyzeGoFileInput,
) (*mcp.CallToolResult, AnalyzeGoFileOutput, error) {
	if strings.TrimSpace(in.Path) == "" {
		return nil, AnalyzeGoFileOutput{}, fmt.Errorf("path cannot be empty")
	}

	analysis, err := internalanalyzer.Analyze(filepath.Clean(in.Path))
	if err != nil {
		return nil, AnalyzeGoFileOutput{}, fmt.Errorf(
			"failed to analyze Go file: %w",
			err,
		)
	}

	return nil, AnalyzeGoFileOutput{
		Imports:   analysis.Imports,
		Functions: analysis.Functions,
		Structs:   analysis.Structs,
	}, nil
}
