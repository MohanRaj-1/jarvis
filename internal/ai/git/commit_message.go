// Package git contains AI-powered Git services.
package git

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"jarvis/internal/ai"
	"jarvis/internal/ai/prompts"
	internalgit "jarvis/internal/git"
)

// CommitMessageService generates Conventional Commit messages from Git diffs.
type CommitMessageService struct {
	Git internalgit.Repository
	AI  ai.Client
}

type commitMessageResponse struct {
	Type    string `json:"type"`
	Scope   string `json:"scope"`
	Subject string `json:"subject"`
}

// GenerateCommitMessage creates a commit message for the uncommitted changes
// in repoPath.
func (s CommitMessageService) GenerateCommitMessage(ctx context.Context, repoPath string) (*CommitMessage, error) {
	if s.Git == nil {
		return nil, errors.New("Git repository is required to generate a commit message")
	}
	if s.AI == nil {
		return nil, errors.New("AI client is required to generate a commit message")
	}

	diff, err := s.Git.Diff(repoPath)
	if err != nil {
		return nil, fmt.Errorf("get Git diff: %w", err)
	}
	if strings.TrimSpace(diff) == "" {
		return nil, errors.New("cannot generate a commit message: repository has no uncommitted changes")
	}

	response, err := s.AI.Generate(ctx, prompts.CommitMessagePrompt(diff))
	if err != nil {
		return nil, fmt.Errorf("generate commit message: %w", err)
	}

	message, err := parseCommitMessage(response)
	if err != nil {
		return nil, fmt.Errorf("parse generated commit message: %w", err)
	}

	return message, nil
}

func parseCommitMessage(response string) (*CommitMessage, error) {
	var message commitMessageResponse
	if err := json.Unmarshal([]byte(response), &message); err != nil {
		return nil, fmt.Errorf("expected JSON with type, scope, and subject: %w", err)
	}

	message.Type = strings.TrimSpace(message.Type)
	message.Scope = strings.TrimSpace(message.Scope)
	message.Subject = strings.TrimSpace(message.Subject)
	if message.Type == "" || message.Subject == "" {
		return nil, errors.New("type and subject are required")
	}
	if strings.ContainsAny(message.Type+message.Scope+message.Subject, "\r\n") {
		return nil, errors.New("commit message fields must be single-line")
	}
	if utf8.RuneCountInString(message.Subject) >= 72 {
		return nil, errors.New("commit message subject must be under 72 characters")
	}

	return &CommitMessage{
		Type:    message.Type,
		Scope:   message.Scope,
		Subject: message.Subject,
	}, nil
}
