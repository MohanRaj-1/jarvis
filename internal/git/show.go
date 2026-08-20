package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/merkletrie"
)

// CommitDetails represents detailed information about a Git commit.
type CommitDetails struct {
	Commit

	Parents []string
	Files   []ChangedFile
}

// ChangedFile describes a file changed by a commit.
type ChangedFile struct {
	Path   string
	Status string // Added, Modified, Deleted, Renamed
}

// Show returns detailed information for hash in repoPath. File changes are
// calculated relative to the first parent, matching Git's default commit view.
func Show(repoPath, hash string) (*CommitDetails, error) {
	if strings.TrimSpace(hash) == "" {
		return nil, fmt.Errorf("commit hash is required")
	}

	repository, cleanPath, err := openRepository(repoPath)
	if err != nil {
		return nil, err
	}

	resolvedHash, err := repository.ResolveRevision(plumbing.Revision(hash))
	if err != nil {
		return nil, fmt.Errorf("resolve commit %q in %q: %w", hash, cleanPath, err)
	}
	commit, err := repository.CommitObject(*resolvedHash)
	if err != nil {
		return nil, fmt.Errorf("read commit %q in %q: %w", hash, cleanPath, err)
	}

	details := &CommitDetails{
		Commit: Commit{
			Hash:    shortHash(commit.Hash),
			Author:  commit.Author.Name,
			Email:   commit.Author.Email,
			Date:    commit.Author.When,
			Message: commit.Message,
		},
		Parents: make([]string, len(commit.ParentHashes)),
	}
	for i, parent := range commit.ParentHashes {
		details.Parents[i] = shortHash(parent)
	}

	files, err := changedFiles(commit)
	if err != nil {
		return nil, fmt.Errorf("get changed files for commit %q in %q: %w", hash, cleanPath, err)
	}
	details.Files = files
	return details, nil
}

func changedFiles(commit *object.Commit) ([]ChangedFile, error) {
	commitTree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	parentTree := &object.Tree{}
	if commit.NumParents() > 0 {
		parent, err := commit.Parent(0)
		if err != nil {
			return nil, err
		}
		parentTree, err = parent.Tree()
		if err != nil {
			return nil, err
		}
	}

	changes, err := parentTree.Diff(commitTree)
	if err != nil {
		return nil, err
	}
	files := make([]ChangedFile, 0, len(changes))
	for _, change := range changes {
		action, err := change.Action()
		if err != nil {
			return nil, err
		}

		file := ChangedFile{}
		switch action {
		case merkletrie.Insert:
			file.Path, file.Status = change.To.Name, "Added"
		case merkletrie.Delete:
			file.Path, file.Status = change.From.Name, "Deleted"
		case merkletrie.Modify:
			file.Path, file.Status = change.To.Name, "Modified"
			if change.From.Name != change.To.Name {
				file.Status = "Renamed"
			}
		default:
			return nil, fmt.Errorf("unsupported file change action %v", action)
		}
		files = append(files, file)
	}
	return files, nil
}

func shortHash(hash plumbing.Hash) string {
	return hash.String()[:7]
}
