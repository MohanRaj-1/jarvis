package analyzer

type Analysis struct {
	Imports   []string
	Functions []string
	Structs   []string
}

func Analyze(path string) (*Analysis, error) {
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

	return &Analysis{
		Imports:   imports,
		Functions: functions,
		Structs:   structs,
	}, nil
}
