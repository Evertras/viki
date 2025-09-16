package viki

import (
	"path/filepath"
	"strings"
)

func filepathToEscapedHttpPath(path string) string {
	// Add other characters to escape as needed, but doing a full URL escape
	// ruins the slashes...
	return strings.ReplaceAll(filepath.ToSlash(path), " ", "%20")
}

func mdPathToHtmlPath(mdPath string) string {
	mdPath = filepath.ToSlash(mdPath)
	return strings.TrimSuffix(mdPath, ".md") + ".html"
}
