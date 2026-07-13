package analyzer

import (
	"go/parser"
	"go/token"
	"strings"
)

type Todo struct {
	Line int    `json:"line"`
	Text string `json:"text"`
}

func ExtractTodos(path string) ([]Todo, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	todos := make([]Todo, 0)
	for _, group := range file.Comments {
		for _, comment := range group.List {
			text := cleanCommentText(comment.Text)
			lines := strings.Split(text, "\n")
			startLine := fileSet.Position(comment.Pos()).Line

			for index, line := range lines {
				todoText, ok := todoText(line)
				if !ok {
					continue
				}

				todos = append(todos, Todo{
					Line: startLine + index,
					Text: todoText,
				})
			}
		}
	}

	return todos, nil
}

func cleanCommentText(text string) string {
	if strings.HasPrefix(text, "//") {
		return strings.TrimSpace(strings.TrimPrefix(text, "//"))
	}

	if strings.HasPrefix(text, "/*") && strings.HasSuffix(text, "*/") {
		text = strings.TrimPrefix(text, "/*")
		text = strings.TrimSuffix(text, "*/")
	}

	return strings.Trim(text, " \t\r")
}

func todoText(line string) (string, bool) {
	line = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "*"))

	const marker = "TODO:"
	before, after, found := strings.Cut(line, marker)
	if !found || strings.TrimSpace(before) != "" {
		return "", false
	}

	return strings.TrimSpace(after), true
}
