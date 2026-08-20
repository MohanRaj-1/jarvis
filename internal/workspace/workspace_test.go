package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewRejectsEmptyRoot(t *testing.T) {
	if _, err := New(" "); err == nil {
		t.Fatal("New() succeeded with an empty root")
	}
}

func TestNewRejectsFileRoot(t *testing.T) {
	file := filepath.Join(t.TempDir(), "file.txt")
	if err := os.WriteFile(file, []byte("test"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := New(file); err == nil {
		t.Fatal("New() succeeded with a file as workspace root")
	}
}

func TestResolve(t *testing.T) {
	root := t.TempDir()
	main := filepath.Join(root, "main.go")
	nested := filepath.Join(root, "internal", "ai", "client.go")
	for _, path := range []string{main, nested} {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("package example"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	w, err := New(root)
	if err != nil {
		t.Fatal(err)
	}
	canonicalRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	canonicalMain, err := filepath.EvalSymlinks(main)
	if err != nil {
		t.Fatal(err)
	}
	canonicalNested, err := filepath.EvalSymlinks(nested)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr string
	}{
		{name: "relative root file", path: "main.go", want: canonicalMain},
		{name: "nested path", path: filepath.Join("internal", "ai", "client.go"), want: canonicalNested},
		{name: "parent traversal outside workspace", path: filepath.Join("..", "secret.txt"), wantErr: "outside workspace"},
		{name: "absolute path outside workspace", path: filepath.Join(filepath.Dir(root), "outside.txt"), wantErr: "outside workspace"},
		{name: "empty path", path: " ", wantErr: "path is required"},
		{name: "workspace root", path: ".", want: canonicalRoot},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := w.Resolve(test.path)
			if test.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantErr) {
					t.Fatalf("Resolve(%q) error = %v, want error containing %q", test.path, err, test.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Resolve(%q) error = %v", test.path, err)
			}
			if got != test.want {
				t.Fatalf("Resolve(%q) = %q, want %q", test.path, got, test.want)
			}
		})
	}
}

func TestResolveRejectsSymlinkEscapingWorkspace(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	link := filepath.Join(root, "secret")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	w, err := New(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Resolve(filepath.Join("secret", "private.txt")); err == nil {
		t.Fatal("Resolve() succeeded for a symlink escaping the workspace")
	}
}
