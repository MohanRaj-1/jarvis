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

const repositorySummaryLogLimit = 10

// RepositorySummary describes the repository's current state.
type RepositorySummary struct {
	Overview        string   `json:"overview"`
	CurrentBranch   string   `json:"current_branch"`
	RecentWork      []string `json:"recent_work"`
	WorkingChanges  []string `json:"working_changes"`
	Recommendations []string `json:"recommendations"`
}

// RepositorySummaryService generates a grounded summary of a Git repository.
type RepositorySummaryService struct {
	Git RepositorySummaryRepository
	AI  ai.Client
}

// RepositorySummaryRepository provides the existing Git operations needed to
// summarize a repository.
type RepositorySummaryRepository interface {
	CurrentBranch(repoPath string) (string, error)
	RepositoryStatus(repoPath string) (*internalgit.Status, error)
	Log(repoPath string, limit int) ([]internalgit.Commit, error)
	Diff(repoPath string) (string, error)
}

// SummarizeRepository gathers the repository state and generates a structured summary.
func (s RepositorySummaryService) SummarizeRepository(ctx context.Context, repoPath string) (*RepositorySummary, error) {
	if s.Git == nil {
		return nil, errors.New("Git repository is required to summarize a repository")
	}
	if s.AI == nil {
		return nil, errors.New("AI client is required to summarize a repository")
	}

	branch, err := s.Git.CurrentBranch(repoPath)
	if err != nil {
		return nil, fmt.Errorf("get current Git branch: %w", err)
	}
	status, err := s.Git.RepositoryStatus(repoPath)
	if err != nil {
		return nil, fmt.Errorf("get Git repository status: %w", err)
	}
	commits, err := s.Git.Log(repoPath, repositorySummaryLogLimit)
	if err != nil {
		return nil, fmt.Errorf("get recent Git commits: %w", err)
	}
	diff, err := s.Git.Diff(repoPath)
	if err != nil {
		return nil, fmt.Errorf("get Git diff: %w", err)
	}

	response, err := s.AI.Generate(ctx, prompts.RepositorySummaryPrompt(branch, status, commits, diff))
	if err != nil {
		return nil, fmt.Errorf("generate repository summary: %w", err)
	}

	summary, err := parseRepositorySummary(response)
	if err != nil {
		return nil, fmt.Errorf("parse generated repository summary: %w", err)
	}
	// The branch is authoritative Git data, rather than AI-generated content.
	summary.CurrentBranch = branch
	return summary, nil
}

func parseRepositorySummary(response string) (*RepositorySummary, error) {
	var summary RepositorySummary
	if err := json.Unmarshal([]byte(response), &summary); err != nil {
		return nil, fmt.Errorf("expected JSON repository summary: %w", err)
	}

	summary.Overview = strings.TrimSpace(summary.Overview)
	if summary.Overview == "" {
		return nil, errors.New("overview is required")
	}
	summary.RecentWork = cleanSummaryEntries(summary.RecentWork)
	summary.WorkingChanges = cleanSummaryEntries(summary.WorkingChanges)
	summary.Recommendations = cleanSummaryEntries(summary.Recommendations)
	return &summary, nil
}

func cleanSummaryEntries(entries []string) []string {
	if entries == nil {
		return []string{}
	}
	for i := range entries {
		entries[i] = strings.TrimSpace(entries[i])
	}
	return entries
}
