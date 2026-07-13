package analyzer

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExtractInterfaces(t *testing.T) {
	source := `package sample

type User struct{}

type Reader interface {
	Read() error
}

type (
	Writer interface {
		Write([]byte) (int, error)
		Close() error
	}
	NotInterface string
	Combined interface {
		Reader
		Writer
		Flush() error
	}
)
`

	path := filepath.Join(t.TempDir(), "sample.go")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	interfaces, err := ExtractInterfaces(path)
	if err != nil {
		t.Fatalf("ExtractInterfaces returned error: %v", err)
	}

	want := []Interface{
		{Name: "Reader", Methods: []string{"Read"}},
		{Name: "Writer", Methods: []string{"Write", "Close"}},
		{Name: "Combined", Methods: []string{"Flush"}},
	}
	if !reflect.DeepEqual(interfaces, want) {
		t.Fatalf("ExtractInterfaces() = %v, want %v", interfaces, want)
	}
}
