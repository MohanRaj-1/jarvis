package git

import (
	"os"
	"path/filepath"
	"testing"

	internalworkspace "jarvis/internal/workspace"
)

func newTestRepository(t *testing.T, root string) DefaultRepository {
	t.Helper()
	workspace, err := internalworkspace.New(root)
	if err != nil {
		t.Fatal(err)
	}
	return NewDefaultRepository(workspace)
}

func TestOperationsRejectRepositoryOutsideWorkspace(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "workspace")
	outside := filepath.Join(parent, "outside")
	for _, path := range []string{root, outside} {
		if err := os.Mkdir(path, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	repository := newTestRepository(t, root)
	tests := []struct {
		name string
		call func() error
	}{
		{name: "current branch", call: func() error { _, err := repository.CurrentBranch(outside); return err }},
		{name: "branches", call: func() error { _, err := repository.RepositoryBranches(outside); return err }},
		{name: "status", call: func() error { _, err := repository.RepositoryStatus(outside); return err }},
		{name: "log", call: func() error { _, err := repository.Log(outside, 1); return err }},
		{name: "diff", call: func() error { _, err := repository.Diff(outside); return err }},
		{name: "show", call: func() error { _, err := repository.Show(outside, "HEAD"); return err }},
		{name: "log range", call: func() error { _, err := repository.LogRange(outside, "HEAD", "HEAD"); return err }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.call(); err == nil {
				t.Fatal("operation succeeded with a repository outside the workspace")
			}
		})
	}
}
