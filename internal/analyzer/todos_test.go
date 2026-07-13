package analyzer

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExtractTodos(t *testing.T) {
	source := `package sample

func Run() {
	// TODO: retry on failure
}

/*
 * TODO: clean up temporary files
 */
`

	path := filepath.Join(t.TempDir(), "sample.go")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	todos, err := ExtractTodos(path)
	if err != nil {
		t.Fatalf("ExtractTodos returned error: %v", err)
	}

	want := []Todo{
		{Line: 4, Text: "retry on failure"},
		{Line: 8, Text: "clean up temporary files"},
	}
	if !reflect.DeepEqual(todos, want) {
		t.Fatalf("ExtractTodos() = %v, want %v", todos, want)
	}
}
