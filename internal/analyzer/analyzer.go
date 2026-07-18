package analyzer

import "fmt"

// Analysis contains structural information extracted from a Go source file.
type Analysis struct {
	Package    string
	Imports    []string
	Functions  []string
	Structs    []string
	Methods    []Method
	Interfaces []Interface
	Todos      []Todo
}

// Analyze extracts all supported structural information from a Go source file.
func Analyze(path string) (*Analysis, error) {
	packageName, err := ExtractPackage(path)
	if err != nil {
		return nil, fmt.Errorf("extract package from %q: %w", path, err)
	}

	imports, err := ExtractImports(path)
	if err != nil {
		return nil, fmt.Errorf("extract imports from %q: %w", path, err)
	}

	functions, err := ExtractFunctions(path)
	if err != nil {
		return nil, fmt.Errorf("extract functions from %q: %w", path, err)
	}

	structs, err := ExtractStructs(path)
	if err != nil {
		return nil, fmt.Errorf("extract structs from %q: %w", path, err)
	}

	methods, err := ExtractMethods(path)
	if err != nil {
		return nil, fmt.Errorf("extract methods from %q: %w", path, err)
	}

	interfaces, err := ExtractInterfaces(path)
	if err != nil {
		return nil, fmt.Errorf("extract interfaces from %q: %w", path, err)
	}

	todos, err := ExtractTodos(path)
	if err != nil {
		return nil, fmt.Errorf("extract TODOs from %q: %w", path, err)
	}

	return &Analysis{
		Package:    packageName,
		Imports:    imports,
		Functions:  functions,
		Structs:    structs,
		Methods:    methods,
		Interfaces: interfaces,
		Todos:      todos,
	}, nil
}
