package analyzer

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExtractImports(t *testing.T) {
	source := `package sample

import (
	"fmt"
	alias "net/http"
	_ "github.com/lib/pq"
)
`

	path := filepath.Join(t.TempDir(), "sample.go")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	imports, err := ExtractImports(path)
	if err != nil {
		t.Fatalf("ExtractImports returned error: %v", err)
	}

	want := []string{"fmt", "net/http", "github.com/lib/pq"}
	if !reflect.DeepEqual(imports, want) {
		t.Fatalf("ExtractImports() = %v, want %v", imports, want)
	}
}
