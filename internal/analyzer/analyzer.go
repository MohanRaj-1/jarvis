package analyzer

type Analysis struct {
	Package    string
	Imports    []string
	Functions  []string
	Structs    []string
	Methods    []Method
	Interfaces []Interface
	Todos      []Todo
}

func Analyze(path string) (*Analysis, error) {
	packageName, err := ExtractPackage(path)
	if err != nil {
		return nil, err
	}

	imports, err := ExtractImports(path)
	if err != nil {
		return nil, err
	}

	functions, err := ExtractFunctions(path)
	if err != nil {
		return nil, err
	}

	structs, err := ExtractStructs(path)
	if err != nil {
		return nil, err
	}

	methods, err := ExtractMethods(path)
	if err != nil {
		return nil, err
	}

	interfaces, err := ExtractInterfaces(path)
	if err != nil {
		return nil, err
	}

	todos, err := ExtractTodos(path)
	if err != nil {
		return nil, err
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
