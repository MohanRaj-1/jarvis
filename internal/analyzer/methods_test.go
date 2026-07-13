package analyzer

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExtractMethods(t *testing.T) {
	source := `package sample

type User struct{}

func NewUser() User {
	return User{}
}

func (u User) Save() {}

func (u *User) Delete() {}
`

	path := filepath.Join(t.TempDir(), "sample.go")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	methods, err := ExtractMethods(path)
	if err != nil {
		t.Fatalf("ExtractMethods returned error: %v", err)
	}

	want := []Method{
		{Receiver: "User", Name: "Save"},
		{Receiver: "*User", Name: "Delete"},
	}
	if !reflect.DeepEqual(methods, want) {
		t.Fatalf("ExtractMethods() = %v, want %v", methods, want)
	}
}
