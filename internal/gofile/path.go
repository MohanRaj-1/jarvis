// Package gofile validates paths to Go source files.
package gofile

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ValidatePath validates and cleans a path to a Go source file.
func ValidatePath(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("path is required; provide a path to a Go source file")
	}

	cleanPath := filepath.Clean(path)
	if !strings.EqualFold(filepath.Ext(cleanPath), ".go") {
		return "", fmt.Errorf("invalid Go file %q; provide a path ending in .go", path)
	}

	return cleanPath, nil
}
