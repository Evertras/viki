package viki

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/afero"
)

func (c *Converter) buildWikiLinkMap(input afero.Fs, inputRootPath string) (map[string]string, error) {
	wikiLinkMap := make(map[string]string)
	inputRootPath = path.Clean(inputRootPath)

	err := afero.Walk(input, inputRootPath, func(inputFilePath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %s: %w", inputFilePath, err)
		}
		if info.IsDir() {
			return nil
		}
		if path.Ext(inputFilePath) != ".md" {
			return nil
		}

		name := strings.TrimSuffix(info.Name(), ".md")

		// If we're running in Windows, convert backslashes to forward slashes
		if os.PathSeparator == '\\' {
			inputFilePath = strings.ReplaceAll(inputFilePath, "\\", "/")
		}
		relativePath := strings.TrimPrefix(inputFilePath, inputRootPath)
		// Make sure to start with a slash
		if !strings.HasPrefix(relativePath, "/") {
			relativePath = "/" + relativePath
		}
		wikiLinkMap[name] = relativePath

		return nil
	})

	return wikiLinkMap, err
}
