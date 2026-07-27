package git

import "testing"

func TestCommitMessageString(t *testing.T) {
	tests := []struct {
		name string
		msg  CommitMessage
		want string
	}{
		{name: "without type", msg: CommitMessage{Subject: "update README"}, want: "update README"},
		{name: "with scope", msg: CommitMessage{Type: "feat", Scope: "ai", Subject: "add commit message generation"}, want: "feat(ai): add commit message generation"},
		{name: "without scope", msg: CommitMessage{Type: "docs", Subject: "update README"}, want: "docs: update README"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}
