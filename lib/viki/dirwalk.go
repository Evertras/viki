package viki

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

func walkDir(input afero.Fs, pathFilter pathFilter, f func(inputFilePath string, info os.FileInfo) error) error {
	return afero.Walk(input, "", func(inputFilePath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path: %w", err)
		}

		if inputFilePath == "." ||
			!pathFilter.isPathIncluded(inputFilePath, info.IsDir()) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		return f(inputFilePath, info)
	})
}
