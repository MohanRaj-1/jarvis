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

func TestReleaseNotesServiceCategorizedResponses(t *testing.T) {
	tests := []struct {
		name         string
		commits      []internalgit.ReleaseCommit
		response     string
		wantFeatures bool
		wantFixes    bool
		wantChanges  bool
	}{
		{
			name: "features",
			commits: []internalgit.ReleaseCommit{
				{Message: "feat(git): add repository status"},
				{Message: "feat(ai): add commit message generation"},
			},
			response:     `{"summary":"Adds repository status and commit message generation.","features":["Added repository status.","Added commit message generation."],"fixes":[],"changes":[],"breaking_changes":[]}`,
			wantFeatures: true,
		},
		{
			name:      "fix",
			commits:   []internalgit.ReleaseCommit{{Message: "fix(ai): handle malformed Gemini response"}},
			response:  `{"summary":"Improves handling of malformed Gemini responses.","features":[],"fixes":["Handle malformed Gemini responses."],"changes":[],"breaking_changes":[]}`,
			wantFixes: true,
		},
		{
			name: "mixed",
			commits: []internalgit.ReleaseCommit{
				{Message: "feat(git): add diff review"},
				{Message: "fix(git): handle empty diff"},
				{Message: "refactor(ai): extract prompt builder"},
			},
			response:     `{"summary":"Adds diff review and improves Git tooling.","features":["Added diff review."],"fixes":["Handle empty diffs."],"changes":["Refactored prompt builder."],"breaking_changes":[]}`,
			wantFeatures: true,
			wantFixes:    true,
			wantChanges:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := git.ReleaseNotesService{
				Git: fakeRepository{commits: tt.commits},
				AI:  &fakeAI{message: tt.response},
			}

			notes, err := service.GenerateReleaseNotes(context.Background(), "/repo", "v0.5.2", "HEAD")
			if err != nil {
				t.Fatalf("GenerateReleaseNotes() error = %v", err)
			}
			if (len(notes.Features) > 0) != tt.wantFeatures {
				t.Errorf("Features = %#v, want non-empty = %t", notes.Features, tt.wantFeatures)
			}
			if (len(notes.Fixes) > 0) != tt.wantFixes {
				t.Errorf("Fixes = %#v, want non-empty = %t", notes.Fixes, tt.wantFixes)
			}
			if (len(notes.Changes) > 0) != tt.wantChanges {
				t.Errorf("Changes = %#v, want non-empty = %t", notes.Changes, tt.wantChanges)
			}
			if len(notes.BreakingChanges) != 0 {
				t.Errorf("BreakingChanges = %#v, want empty", notes.BreakingChanges)
			}
		})
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
