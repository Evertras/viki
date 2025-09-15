package viki

import (
	"fmt"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/spf13/afero"
)

type pathFilter struct {
	ignoreChecker  *ignore.GitIgnore
	includeChecker *ignore.GitIgnore
}

func (p *pathFilter) isPathIncluded(path string, isDir bool) bool {
	if isDir {
		// For directories, we only care about the ignore patterns
		return !p.ignoreChecker.MatchesPath(path)
	}

	return !p.ignoreChecker.MatchesPath(path) && p.includeChecker.MatchesPath(path)
}

func generatePathFilter(cfg ConverterOptions, input afero.Fs) (pathFilter, error) {
	const gitignoreFilename = ".gitignore"

	ignoreLines := make([]string, len(cfg.ExcludePatterns))
	copy(ignoreLines, cfg.ExcludePatterns)

	exists, err := afero.Exists(input, gitignoreFilename)
	if err != nil {
		return pathFilter{}, fmt.Errorf("failed to check for .gitignore: %w", err)
	}

	if exists {
		content, err := afero.ReadFile(input, gitignoreFilename)
		if err != nil {
			return pathFilter{}, fmt.Errorf("failed to read .gitignore: %w", err)
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

	ignoreChecker := ignore.CompileIgnoreLines(cleanedIgnoreLines...)

	// Avoid modifying the original config slices
	includePatterns := make([]string, len(cfg.IncludePatterns))
	copy(includePatterns, cfg.IncludePatterns)

	if len(includePatterns) == 0 {
		includePatterns = []string{"**"}
	}

	includeChecker := ignore.CompileIgnoreLines(includePatterns...)

	return pathFilter{
		ignoreChecker:  ignoreChecker,
		includeChecker: includeChecker,
	}, nil
}
