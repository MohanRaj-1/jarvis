package git_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"jarvis/internal/ai/git"
	internalgit "jarvis/internal/git"
)

func TestReleaseNotesServiceGenerateReleaseNotes(t *testing.T) {
	ai := &fakeAI{message: `{
		"summary":"Adds Git intelligence features.",
		"features":["Added release note generation."],
		"fixes":[],
		"changes":["Improved Git tooling."]
	}`}
	service := git.ReleaseNotesService{
		Git: fakeRepository{commits: []internalgit.ReleaseCommit{{Hash: "abcdef0", Author: "Mohan Raj", Message: "feat(git): add release notes", ChangedFiles: []internalgit.ChangedFile{{Path: "internal/ai/git/release_notes.go", Status: "Added"}}}}},
		AI:  ai,
	}

	notes, err := service.GenerateReleaseNotes(context.Background(), "/repo", "v0.5.2", "HEAD")
	if err != nil {
		t.Fatalf("GenerateReleaseNotes() error = %v", err)
	}
	if notes.Summary != "Adds Git intelligence features." {
		t.Errorf("GenerateReleaseNotes().Summary = %q", notes.Summary)
	}
	if len(notes.Features) != 1 || notes.Features[0] != "Added release note generation." {
		t.Errorf("GenerateReleaseNotes().Features = %#v", notes.Features)
	}
	for _, want := range []string{"v0.5.2", "HEAD", "abcdef0", "Mohan Raj", "feat(git): add release notes", "internal/ai/git/release_notes.go", "Added"} {
		if !strings.Contains(ai.prompt, want) {
			t.Errorf("prompt does not contain %q: %q", want, ai.prompt)
		}
	}
}

func TestReleaseNotesServiceGenerateReleaseNotesErrors(t *testing.T) {
	tests := []struct {
		name    string
		service git.ReleaseNotesService
		want    string
	}{
		{name: "missing Git repository", service: git.ReleaseNotesService{AI: &fakeAI{}}, want: "Git repository is required"},
		{name: "missing AI client", service: git.ReleaseNotesService{Git: fakeRepository{}}, want: "AI client is required"},
		{name: "empty range", service: git.ReleaseNotesService{Git: fakeRepository{}, AI: &fakeAI{}}, want: "range contains no commits"},
		{name: "Git failure", service: git.ReleaseNotesService{Git: fakeRepository{logErr: errors.New("failed")}, AI: &fakeAI{}}, want: "get Git commit range"},
		{name: "AI failure", service: git.ReleaseNotesService{Git: fakeRepository{commits: []internalgit.ReleaseCommit{{Message: "feat: add feature"}}}, AI: &fakeAI{err: errors.New("failed")}}, want: "generate release notes"},
		{name: "invalid JSON", service: git.ReleaseNotesService{Git: fakeRepository{commits: []internalgit.ReleaseCommit{{Message: "feat: add feature"}}}, AI: &fakeAI{message: "not JSON"}}, want: "parse generated release notes"},
		{name: "missing summary", service: git.ReleaseNotesService{Git: fakeRepository{commits: []internalgit.ReleaseCommit{{Message: "feat: add feature"}}}, AI: &fakeAI{message: `{"features":[]}`}}, want: "summary is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.service.GenerateReleaseNotes(context.Background(), "/repo", "v1.0.0", "HEAD")
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("GenerateReleaseNotes() error = %v, want %q", err, tt.want)
			}
		})
	}
}
