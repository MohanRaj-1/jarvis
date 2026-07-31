package git

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"jarvis/internal/ai"
	"jarvis/internal/ai/prompts"
	internalgit "jarvis/internal/git"
)

// ExplainCommitService generates grounded, human-readable commit explanations.
type ExplainCommitService struct {
	Git internalgit.Repository
	AI  ai.Client
}

// ExplainCommit explains the commit identified by hash in repoPath.
func (s ExplainCommitService) ExplainCommit(ctx context.Context, repoPath, hash string) (string, error) {
	if s.Git == nil {
		return "", errors.New("Git repository is required to explain a commit")
	}
	if s.AI == nil {
		return "", errors.New("AI client is required to explain a commit")
	}

	commit, err := s.Git.Show(repoPath, hash)
	if err != nil {
		return "", fmt.Errorf("show Git commit: %w", err)
	}

	explanation, err := s.AI.Generate(ctx, prompts.ExplainCommitPrompt(commit))
	if err != nil {
		return "", fmt.Errorf("generate commit explanation: %w", err)
	}

	explanation = strings.TrimSpace(explanation)
	if explanation == "" {
		return "", errors.New("generated commit explanation is empty")
	}

	return explanation, nil
}
