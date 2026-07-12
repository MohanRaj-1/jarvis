package analyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func ExtractFunctions(path string) ([]string, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, 0)
	if err != nil {
		return nil, err
	}

	functions := make([]string, 0)
	ast.Inspect(file, func(node ast.Node) bool {
		functionDeclaration, ok := node.(*ast.FuncDecl)
		if !ok {
			return true
		}

		functions = append(functions, functionDeclaration.Name.Name)
		// Skip traversing inside the function body since
		// nested function declarations are not allowed in Go.
		return false
	})

	return functions, nil
}
