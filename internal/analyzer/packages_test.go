package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractPackage(t *testing.T) {
	source := `package workspace

func Run() {}
`

	path := filepath.Join(t.TempDir(), "sample.go")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	packageName, err := ExtractPackage(path)
	if err != nil {
		t.Fatalf("ExtractPackage returned error: %v", err)
	}

	if packageName != "workspace" {
		t.Fatalf("ExtractPackage() = %q, want %q", packageName, "workspace")
	}
}
