package analyzer

import (
	"go/parser"
	"go/token"
)

func ExtractPackage(path string) (string, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, parser.PackageClauseOnly)
	if err != nil {
		return "", err
	}

	return file.Name.Name, nil
}
