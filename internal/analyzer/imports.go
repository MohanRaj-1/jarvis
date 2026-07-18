package analyzer

import (
	"fmt"
	"go/parser"
	"go/token"
	"strconv"
)

// ExtractImports returns import paths declared in a Go source file.
func ExtractImports(path string) ([]string, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, parser.ImportsOnly)
	if err != nil {
		return nil, fmt.Errorf("parse imports in Go source file %q: %w", path, err)
	}

	imports := make([]string, 0, len(file.Imports))
	for _, spec := range file.Imports {
		importPath, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			return nil, fmt.Errorf("decode import path in %q: %w", path, err)
		}

		imports = append(imports, importPath)
	}

	return imports, nil
}
