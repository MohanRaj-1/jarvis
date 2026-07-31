// Package prompts contains prompt builders for AI features.
package prompts

import (
	"fmt"
	"strings"

	internalgit "jarvis/internal/git"
)

// CommitMessagePrompt creates the prompt used to generate structured commit
// message data from a Git diff.
func CommitMessagePrompt(diff string) string {
	return fmt.Sprintf(`You are an expert software engineer.

Analyze the Git diff and generate the most important Conventional Commit.

Rules:
- Return valid JSON only. Do not use Markdown or code fences.
- Do not explain your reasoning.
- The JSON must have exactly these string fields: "type", "scope", and "subject".
- "type" must be a concise, lowercase Conventional Commit value. "scope" may be an empty string.
- "subject" must be one line and under 72 characters.

Example response:
{"type":"feat","scope":"ai","subject":"add AI-powered commit message generation"}

Git diff:
%s`, diff)
}

// ExplainCommitPrompt creates a readable, grounded prompt from commit details.
func ExplainCommitPrompt(commit *internalgit.CommitDetails) string {
	var files strings.Builder
	if len(commit.Files) == 0 {
		files.WriteString("- No file changes recorded")
	} else {
		for _, file := range commit.Files {
			fmt.Fprintf(&files, "- %s (%s)\n", file.Path, file.Status)
		}
	}

	return fmt.Sprintf(`You are a senior software engineer.

Explain the following Git commit.

Rules:
- Use plain English.
- Maximum 150 words.
- Do not repeat the commit message.
- Focus on what changed and why it matters.
- Do not speculate beyond the provided information.

Commit Message:
%s

Author:
%s

Files:
%s`, strings.TrimSpace(commit.Message), strings.TrimSpace(commit.Author), strings.TrimSpace(files.String()))
}
