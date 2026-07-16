package ai

import (
	"encoding/json"
	"fmt"

	"jarvis/internal/analyzer"
)

func BuildExplainPrompt(analysis *analyzer.Analysis, source string) string {
	analysisJSON, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		analysisJSON = []byte("null")
	}

	return fmt.Sprintf(`You are an expert Go developer. Explain the following Go source code clearly and concisely.

Use the static analysis as supporting context, but rely on the source code as the authoritative reference. Cover the code's purpose, its important types and functions, how its dependencies are used, and any TODO items or notable concerns.

Static analysis:
%s

Source code:
~~~go
%s
~~~
`, analysisJSON, source)
}
