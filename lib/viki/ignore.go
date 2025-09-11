package viki

import (
	"fmt"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/spf13/afero"
)

func (c *Converter) generateIgnoreChecker(fs afero.Fs, basePath string) (*ignore.GitIgnore, error) {
	var ignoreChecker *ignore.GitIgnore

	ignoreLines := make([]string, len(c.config.ExcludePatterns))
	copy(ignoreLines, c.config.ExcludePatterns)

	exists, err := afero.Exists(fs, filepath.Join(basePath, ".gitignore"))
	if err != nil {
		return nil, fmt.Errorf("failed to check for .gitignore: %w", err)
	}

	if exists {
		content, err := afero.ReadFile(fs, filepath.Join(basePath, ".gitignore"))
		if err != nil {
			return nil, fmt.Errorf("failed to read .gitignore: %w", err)
		}

		ignoreLines = append(ignoreLines, strings.Split(string(content), "\n")...)
	}

	// filter out all ignorelines that are empty or just whitespace
	cleanedIgnoreLines := []string{}
	for _, line := range ignoreLines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanedIgnoreLines = append(cleanedIgnoreLines, line)
		}
	}
	ignoreLines = cleanedIgnoreLines

	ignoreChecker = ignore.CompileIgnoreLines(ignoreLines...)

	return ignoreChecker, nil
}
