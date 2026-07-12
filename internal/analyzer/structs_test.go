package analyzer

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExtractStructs(t *testing.T) {
	source := `package sample

type Single struct {
	Name string
}

type Alias = Single

type (
	Grouped struct {
		ID int
	}
	NotStruct string
)
`

	path := filepath.Join(t.TempDir(), "sample.go")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	structs, err := ExtractStructs(path)
	if err != nil {
		t.Fatalf("ExtractStructs returned error: %v", err)
	}

	want := []string{"Single", "Grouped"}
	if !reflect.DeepEqual(structs, want) {
		t.Fatalf("ExtractStructs() = %v, want %v", structs, want)
	}
}
