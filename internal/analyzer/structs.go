package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// ExtractStructs returns structs declared in a Go source file.
func ExtractStructs(path string) ([]string, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("parse Go source file %q: %w", path, err)
	}

	structs := make([]string, 0)
	ast.Inspect(file, func(node ast.Node) bool {
		generalDeclaration, ok := node.(*ast.GenDecl)
		if !ok || generalDeclaration.Tok != token.TYPE {
			return true
		}

		for _, spec := range generalDeclaration.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			if _, ok := typeSpec.Type.(*ast.StructType); ok {
				structs = append(structs, typeSpec.Name.Name)
			}
		}

		return false
	})

	return structs, nil
}
