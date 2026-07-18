package analyzer

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestAnalyze(t *testing.T) {
	source := `package sample

import "fmt"

// TODO: replace the placeholder output
type User struct{}

type Saver interface {
	Save() error
}

func NewUser() User { return User{} }

func (User) Save() error {
	fmt.Println("saved")
	return nil
}
`

	path := filepath.Join(t.TempDir(), "sample.go")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	analysis, err := Analyze(path)
	if err != nil {
		t.Fatalf("Analyze returned error: %v", err)
	}

	want := &Analysis{
		Package:    "sample",
		Imports:    []string{"fmt"},
		Functions:  []string{"NewUser", "Save"},
		Structs:    []string{"User"},
		Methods:    []Method{{Receiver: "User", Name: "Save"}},
		Interfaces: []Interface{{Name: "Saver", Methods: []string{"Save"}}},
		Todos:      []Todo{{Line: 5, Text: "replace the placeholder output"}},
	}
	if !reflect.DeepEqual(analysis, want) {
		t.Fatalf("Analyze() = %#v, want %#v", analysis, want)
	}
}
