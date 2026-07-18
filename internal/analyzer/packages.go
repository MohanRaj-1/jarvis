package analyzer

import (
	"fmt"
	"go/parser"
	"go/token"
)

// ExtractPackage returns the package name declared in a Go source file.
func ExtractPackage(path string) (string, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, parser.PackageClauseOnly)
	if err != nil {
		return "", fmt.Errorf("parse package declaration in Go source file %q: %w", path, err)
	}

	return file.Name.Name, nil
}
