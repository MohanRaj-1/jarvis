package ai

import (
	"encoding/json"
	"fmt"

	"jarvis/internal/analyzer"
)

// BuildExplainPrompt creates the structured prompt used to explain Go source files.
func BuildExplainPrompt(analysis *analyzer.Analysis, source string) string {
	analysisJSON, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		analysisJSON = []byte("null")
	}

	return fmt.Sprintf(`You are an expert Go developer. Explain the following Go source code clearly and concisely.

Use the static analysis as supporting context, but treat the source code as the authoritative reference. Do not invent or infer behavior that is not supported by the source.

Use exactly these Markdown headings, in this order:

# Purpose
# High-Level Flow
# Important Types
# Important Functions
# External Dependencies
# Interesting Design Decisions
# Possible Improvements

Under each heading, explain why the code is written this way, not only what it does. Quote every function name exactly as it appears in the source code. Mention every TODO comment you find, including its location and intent. If a section has no applicable details, say so briefly rather than speculating.

Static analysis:
%s

Source code:
~~~go
%s
~~~
`, analysisJSON, source)
}
