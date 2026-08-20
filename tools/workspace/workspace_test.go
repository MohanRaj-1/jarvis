package workspace

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	internalworkspace "jarvis/internal/workspace"
)

func TestReadFileReadsPathInsideWorkspace(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "allowed.txt"), []byte("allowed content"), 0o644); err != nil {
		t.Fatal(err)
	}

	ws, err := internalworkspace.New(root)
	if err != nil {
		t.Fatal(err)
	}
	_, output, err := NewReadFileTool(ws).Handle(context.Background(), nil, ReadFileInput{Path: "allowed.txt"})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if output.Content != "allowed content" {
		t.Fatalf("Handle() content = %q, want %q", output.Content, "allowed content")
	}
	if output.Size != int64(len("allowed content")) {
		t.Fatalf("Handle() size = %d, want %d", output.Size, len("allowed content"))
	}
}

func TestToolsRejectPathsOutsideWorkspace(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "workspace")
	outside := filepath.Join(parent, "outside")
	for _, dir := range []string{root, outside} {
		if err := os.Mkdir(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("secret"), 0o644); err != nil {
		t.Fatal(err)
	}

	ws, err := internalworkspace.New(root)
	if err != nil {
		t.Fatal(err)
	}
	outsideFile := filepath.Join("..", "outside", "secret.txt")

	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "read file",
			call: func() error {
				_, _, err := NewReadFileTool(ws).Handle(context.Background(), nil, ReadFileInput{Path: outsideFile})
				return err
			},
		},
		{
			name: "list directory",
			call: func() error {
				_, _, err := NewListDirectoryTool(ws).Handle(context.Background(), nil, ListDirectoryInput{Path: filepath.Join("..", "outside")})
				return err
			},
		},
		{
			name: "file info",
			call: func() error {
				_, _, err := NewFileInfoTool(ws).Handle(context.Background(), nil, FileInfoInput{Path: outsideFile})
				return err
			},
		},
		{
			name: "find files",
			call: func() error {
				_, _, err := NewFindFilesTool(ws).Handle(context.Background(), nil, FindFilesInput{Root: filepath.Join("..", "outside"), Pattern: "*.txt"})
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
