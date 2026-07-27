// Package prompts contains prompt builders for AI features.
package prompts

import "fmt"

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
