package analyzer

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExtractFunctions(t *testing.T) {
	source := `package sample

func Exported() {}

func unexported(value int) int {
	callback := func() {}
	callback()
	return value
}

func (User) Save() {}
`

	path := filepath.Join(t.TempDir(), "sample.go")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	functions, err := ExtractFunctions(path)
	if err != nil {
		t.Fatalf("ExtractFunctions returned error: %v", err)
	}

	want := []string{"Exported", "unexported", "Save"}
	if !reflect.DeepEqual(functions, want) {
		t.Fatalf("ExtractFunctions() = %v, want %v", functions, want)
	}
}
