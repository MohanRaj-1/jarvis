package git

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"jarvis/internal/ai"
	"jarvis/internal/ai/prompts"
	internalgit "jarvis/internal/git"
)

// ReleaseNotes contains structured notes for a Git commit range.
type ReleaseNotes struct {
	Summary         string   `json:"summary"`
	Features        []string `json:"features"`
	Fixes           []string `json:"fixes"`
	Changes         []string `json:"changes"`
	BreakingChanges []string `json:"breaking_changes"`
}

// ReleaseNotesService generates release notes from a Git commit range.
type ReleaseNotesService struct {
	Git internalgit.Repository
	AI  ai.Client
}

// GenerateReleaseNotes creates structured release notes for commits in from..to.
func (s ReleaseNotesService) GenerateReleaseNotes(ctx context.Context, repoPath, from, to string) (*ReleaseNotes, error) {
	if s.Git == nil {
		return nil, errors.New("Git repository is required to generate release notes")
	}
	if s.AI == nil {
		return nil, errors.New("AI client is required to generate release notes")
	}

	commits, err := s.Git.LogRange(repoPath, from, to)
	if err != nil {
		return nil, fmt.Errorf("get Git commit range: %w", err)
	}
	if len(commits) == 0 {
		return nil, errors.New("cannot generate release notes: commit range contains no commits")
	}

	response, err := s.AI.Generate(ctx, prompts.ReleaseNotesPrompt(from, to, commits))
	if err != nil {
		return nil, fmt.Errorf("generate release notes: %w", err)
	}

	notes, err := parseReleaseNotes(response)
	if err != nil {
		return nil, fmt.Errorf("parse generated release notes: %w", err)
	}
	return notes, nil
}

func parseReleaseNotes(response string) (*ReleaseNotes, error) {
	var notes ReleaseNotes
	if err := json.Unmarshal([]byte(response), &notes); err != nil {
		return nil, fmt.Errorf("expected JSON release notes: %w", err)
	}
	notes.Summary = strings.TrimSpace(notes.Summary)
	if notes.Summary == "" {
		return nil, errors.New("summary is required")
	}
	notes.Features = cleanEntries(notes.Features)
	notes.Fixes = cleanEntries(notes.Fixes)
	notes.Changes = cleanEntries(notes.Changes)
	notes.BreakingChanges = cleanEntries(notes.BreakingChanges)
	return &notes, nil
}

func cleanEntries(entries []string) []string {
	for i := range entries {
		entries[i] = strings.TrimSpace(entries[i])
	}
	return entries
}
