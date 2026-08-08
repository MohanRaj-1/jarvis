package git

import (
	"errors"
	"fmt"
	"strings"
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

// ReleaseCommit is the Git history needed to generate grounded release notes.
type ReleaseCommit struct {
	Hash         string
	Author       string
	Date         time.Time
	Message      string
	ChangedFiles []ChangedFile
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

// LogRange returns commits reachable from to and not reachable from from,
// equivalent to Git's "from..to" revision range.
func LogRange(repoPath, from, to string) ([]ReleaseCommit, error) {
	if strings.TrimSpace(from) == "" || strings.TrimSpace(to) == "" {
		return nil, fmt.Errorf("both range endpoints are required")
	}

	repository, cleanPath, err := openRepository(repoPath)
	if err != nil {
		return nil, err
	}

	fromHash, err := repository.ResolveRevision(plumbing.Revision(from))
	if err != nil {
		return nil, fmt.Errorf("resolve range start %q in %q: %w", from, cleanPath, err)
	}
	toHash, err := repository.ResolveRevision(plumbing.Revision(to))
	if err != nil {
		return nil, fmt.Errorf("resolve range end %q in %q: %w", to, cleanPath, err)
	}

	excluded, err := reachableCommitHashes(repository, *fromHash)
	if err != nil {
		return nil, fmt.Errorf("read commits from range start %q in %q: %w", from, cleanPath, err)
	}

	iterator, err := repository.Log(&gitlib.LogOptions{From: *toHash})
	if err != nil {
		return nil, fmt.Errorf("read commits from range end %q in %q: %w", to, cleanPath, err)
	}
	defer iterator.Close()

	commits := []ReleaseCommit{}
	err = iterator.ForEach(func(commit *object.Commit) error {
		if _, skip := excluded[commit.Hash]; skip {
			return nil
		}
		files, err := changedFiles(commit)
		if err != nil {
			return err
		}
		commits = append(commits, ReleaseCommit{
			Hash:         commit.Hash.String()[:7],
			Author:       commit.Author.Name,
			Date:         commit.Author.When,
			Message:      commit.Message,
			ChangedFiles: files,
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("iterate commit range %q..%q in %q: %w", from, to, cleanPath, err)
	}

	return commits, nil
}

func reachableCommitHashes(repository *gitlib.Repository, from plumbing.Hash) (map[plumbing.Hash]struct{}, error) {
	iterator, err := repository.Log(&gitlib.LogOptions{From: from})
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	hashes := make(map[plumbing.Hash]struct{})
	err = iterator.ForEach(func(commit *object.Commit) error {
		hashes[commit.Hash] = struct{}{}
		return nil
	})
	return hashes, err
}
