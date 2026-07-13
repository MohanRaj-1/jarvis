package analyzer

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
)

type Method struct {
	Receiver string `json:"receiver"`
	Name     string `json:"name"`
}

func ExtractMethods(path string) ([]Method, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, 0)
	if err != nil {
		return nil, err
	}

	methods := make([]Method, 0)
	var extractErr error
	ast.Inspect(file, func(node ast.Node) bool {
		if extractErr != nil {
			return false
		}

		functionDeclaration, ok := node.(*ast.FuncDecl)
		if !ok {
			return true
		}

		if functionDeclaration.Recv == nil || len(functionDeclaration.Recv.List) == 0 {
			return false
		}

		receiver, err := receiverType(fileSet, functionDeclaration.Recv.List[0].Type)
		if err != nil {
			extractErr = err
			return false
		}

		methods = append(methods, Method{
			Receiver: receiver,
			Name:     functionDeclaration.Name.Name,
		})

		return false
	})
	if extractErr != nil {
		return nil, extractErr
	}

	return methods, nil
}

func receiverType(fileSet *token.FileSet, expression ast.Expr) (string, error) {
	var buffer bytes.Buffer
	if err := printer.Fprint(&buffer, fileSet, expression); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
