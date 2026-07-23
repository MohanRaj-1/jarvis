package git

import (
	"errors"
	"fmt"
	"time"

	gitlib "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

const defaultLogLimit = 10

// Commit represents a Git commit.
type Commit struct {
	Hash    string
	Author  string
	Email   string
	Date    time.Time
	Message string
}

// Log returns the most recent commits in repoPath. A limit of zero uses the
// default limit of 10.
func Log(repoPath string, limit int) ([]Commit, error) {
	if limit < 0 {
		return nil, fmt.Errorf("commit log limit cannot be negative")
	}
	if limit == 0 {
		limit = defaultLogLimit
	}

	repository, cleanPath, err := openRepository(repoPath)
	if err != nil {
		return nil, err
	}

	iterator, err := repository.Log(&gitlib.LogOptions{})
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			return []Commit{}, nil
		}
		return nil, fmt.Errorf("read commit log for %q: %w", cleanPath, err)
	}
	defer iterator.Close()

	commits := make([]Commit, 0, limit)
	err = iterator.ForEach(func(commit *object.Commit) error {
		if len(commits) == limit {
			return storer.ErrStop
		}
		commits = append(commits, Commit{
			Hash:    commit.Hash.String()[:7],
			Author:  commit.Author.Name,
			Email:   commit.Author.Email,
			Date:    commit.Author.When,
			Message: commit.Message,
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("iterate commit log for %q: %w", cleanPath, err)
	}

	return commits, nil
}
