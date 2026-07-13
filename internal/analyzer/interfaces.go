package analyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type Interface struct {
	Name    string   `json:"name"`
	Methods []string `json:"methods"`
}

func ExtractInterfaces(path string) ([]Interface, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, 0)
	if err != nil {
		return nil, err
	}

	interfaces := make([]Interface, 0)
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

			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}

			interfaces = append(interfaces, Interface{
				Name:    typeSpec.Name.Name,
				Methods: interfaceMethods(interfaceType),
			})
		}

		return false
	})

	return interfaces, nil
}

func interfaceMethods(interfaceType *ast.InterfaceType) []string {
	methods := make([]string, 0)
	if interfaceType.Methods == nil {
		return methods
	}

	for _, field := range interfaceType.Methods.List {
		if _, ok := field.Type.(*ast.FuncType); !ok {
			continue
		}

		for _, name := range field.Names {
			methods = append(methods, name.Name)
		}
	}

	return methods
}
