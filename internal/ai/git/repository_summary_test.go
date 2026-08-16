package git_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"jarvis/internal/ai/git"
	internalgit "jarvis/internal/git"
)

func TestRepositorySummaryServiceSummarizeRepository(t *testing.T) {
	ai := &fakeAI{message: `{
		"overview":"Repository work is in progress.",
		"current_branch":"wrong-branch",
		"recent_work":["Added repository status support."],
		"working_changes":["Modified the Git service."],
		"recommendations":[]
	}`}
	service := git.RepositorySummaryService{
		Git: summaryRepository{
			branch:  "feature/repository-summary",
			status:  &internalgit.Status{Modified: []string{"internal/ai/git/repository_summary.go"}},
			commits: []internalgit.Commit{{Hash: "abcdef0", Message: "feat(git): add repository summary", Date: time.Date(2026, 8, 16, 0, 0, 0, 0, time.UTC)}},
			diff:    "diff --git a/summary.go b/summary.go\n+type RepositorySummary struct{}",
		},
		AI: ai,
	}

	summary, err := service.SummarizeRepository(context.Background(), "/repo")
	if err != nil {
		t.Fatalf("SummarizeRepository() error = %v", err)
	}
	if summary.CurrentBranch != "feature/repository-summary" {
		t.Errorf("CurrentBranch = %q, want gathered branch", summary.CurrentBranch)
	}
	if len(summary.Recommendations) != 0 {
		t.Errorf("Recommendations = %#v, want empty", summary.Recommendations)
	}
	for _, want := range []string{"feature/repository-summary", "internal/ai/git/repository_summary.go", "abcdef0", "feat(git): add repository summary", "diff --git a/summary.go b/summary.go"} {
		if !strings.Contains(ai.prompt, want) {
			t.Errorf("prompt does not contain %q: %q", want, ai.prompt)
		}
	}
}

func TestRepositorySummaryServiceNormalizesMissingArrays(t *testing.T) {
	service := git.RepositorySummaryService{
		Git: summaryRepository{branch: "main", status: &internalgit.Status{}},
		AI:  &fakeAI{message: `{"overview":"The repository is clean."}`},
	}

	summary, err := service.SummarizeRepository(context.Background(), "/repo")
	if err != nil {
		t.Fatalf("SummarizeRepository() error = %v", err)
	}
	if summary.RecentWork == nil || summary.WorkingChanges == nil || summary.Recommendations == nil {
		t.Errorf("summary arrays must be non-nil: %#v", summary)
	}
}

func TestRepositorySummaryServiceErrors(t *testing.T) {
	tests := []struct {
		name    string
		service git.RepositorySummaryService
		want    string
	}{
		{name: "missing Git repository", service: git.RepositorySummaryService{AI: &fakeAI{}}, want: "Git repository is required"},
		{name: "missing AI client", service: git.RepositorySummaryService{Git: summaryRepository{}}, want: "AI client is required"},
		{name: "branch failure", service: git.RepositorySummaryService{Git: summaryRepository{branchErr: errors.New("failed")}, AI: &fakeAI{}}, want: "get current Git branch"},
		{name: "status failure", service: git.RepositorySummaryService{Git: summaryRepository{branch: "main", statusErr: errors.New("failed")}, AI: &fakeAI{}}, want: "get Git repository status"},
		{name: "log failure", service: git.RepositorySummaryService{Git: summaryRepository{branch: "main", status: &internalgit.Status{}, logErr: errors.New("failed")}, AI: &fakeAI{}}, want: "get recent Git commits"},
		{name: "diff failure", service: git.RepositorySummaryService{Git: summaryRepository{branch: "main", status: &internalgit.Status{}, diffErr: errors.New("failed")}, AI: &fakeAI{}}, want: "get Git diff"},
		{name: "AI failure", service: git.RepositorySummaryService{Git: summaryRepository{branch: "main", status: &internalgit.Status{}}, AI: &fakeAI{err: errors.New("failed")}}, want: "generate repository summary"},
		{name: "invalid JSON", service: git.RepositorySummaryService{Git: summaryRepository{branch: "main", status: &internalgit.Status{}}, AI: &fakeAI{message: "not JSON"}}, want: "parse generated repository summary"},
		{name: "missing overview", service: git.RepositorySummaryService{Git: summaryRepository{branch: "main", status: &internalgit.Status{}}, AI: &fakeAI{message: `{"recommendations":[]}`}}, want: "overview is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.service.SummarizeRepository(context.Background(), "/repo")
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("SummarizeRepository() error = %v, want %q", err, tt.want)
			}
		})
	}
}

type summaryRepository struct {
	branch    string
	branchErr error
	status    *internalgit.Status
	statusErr error
	commits   []internalgit.Commit
	logErr    error
	diff      string
	diffErr   error
}

func (r summaryRepository) CurrentBranch(string) (string, error) { return r.branch, r.branchErr }
func (r summaryRepository) RepositoryStatus(string) (*internalgit.Status, error) {
	return r.status, r.statusErr
}
func (r summaryRepository) Log(string, int) ([]internalgit.Commit, error) { return r.commits, r.logErr }
func (r summaryRepository) Diff(string) (string, error)                   { return r.diff, r.diffErr }
