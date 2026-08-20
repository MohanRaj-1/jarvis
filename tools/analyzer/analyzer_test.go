package analyzer

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	internalworkspace "jarvis/internal/workspace"
)

func TestAnalyzeGoFileAnalyzesPathInsideWorkspace(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	ws, err := internalworkspace.New(root)
	if err != nil {
		t.Fatal(err)
	}
	_, output, err := NewAnalyzeGoFileTool(ws).Handle(context.Background(), nil, AnalyzeGoFileInput{Path: "main.go"})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if output.Package != "main" {
		t.Fatalf("Handle() package = %q, want %q", output.Package, "main")
	}
}

func TestAnalyzeImportsAnalyzesPathInsideWorkspace(t *testing.T) {
	root := t.TempDir()
	source := `package main

import "fmt"

func main() {
	fmt.Println("hello")
}
`
	if err := os.WriteFile(filepath.Join(root, "main.go"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}

	ws, err := internalworkspace.New(root)
	if err != nil {
		t.Fatal(err)
	}
	_, output, err := NewAnalyzeImportsTool(ws).Handle(
		context.Background(),
		nil,
		AnalyzeImportsInput{Path: "main.go"},
	)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if len(output.Imports) != 1 || output.Imports[0] != "fmt" {
		t.Fatalf("Handle() imports = %v, want [fmt]", output.Imports)
	}
}

func TestAnalyzerToolsRejectPathsOutsideWorkspace(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "workspace")
	outside := filepath.Join(parent, "outside")
	for _, dir := range []string{root, outside} {
		if err := os.Mkdir(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(outside, "secret.go"), []byte("package secret\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	ws, err := internalworkspace.New(root)
	if err != nil {
		t.Fatal(err)
	}
	outsideFile := filepath.Join("..", "outside", "secret.go")

	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "analyze Go file",
			call: func() error {
				_, _, err := NewAnalyzeGoFileTool(ws).Handle(context.Background(), nil, AnalyzeGoFileInput{Path: outsideFile})
				return err
			},
		},
		{
			name: "analyze imports",
			call: func() error {
				_, _, err := NewAnalyzeImportsTool(ws).Handle(context.Background(), nil, AnalyzeImportsInput{Path: outsideFile})
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.call(); err == nil {
				t.Fatal("tool succeeded with a path outside the workspace")
			}
		})
	}
}
