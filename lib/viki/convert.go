package viki

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/afero"
)

type ConverterOptions struct {
}

type Converter struct {
}

func NewConverter(options ConverterOptions) *Converter {
	return &Converter{}
}

// Convert takes in a filesystem (which could be an OS filesystem or an in-memory one)
// and an output afero.Fs, converts each .md file to .html, and then writes to the output afero.Fs.
func (c *Converter) Convert(input afero.Fs, inputRootPath string, output afero.Fs) error {
	_, err := c.buildWikiLinkMap(input, inputRootPath)
	if err != nil {
		return fmt.Errorf("failed to build wiki link map: %w", err)
	}

	return afero.Walk(input, inputRootPath, func(inputFilePath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %s: %w", inputFilePath, err)
		}
		if info.IsDir() {
			return nil
		}

		if path.Ext(inputFilePath) != ".md" {
			return nil
		}

		content, err := afero.ReadFile(input, inputFilePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", inputFilePath, err)
		}

		// TODO: Convert wikilinks

		outputFilePath := mdPathToHTMLPath(inputFilePath)

		err = afero.WriteFile(output, outputFilePath, content, 0644)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", outputFilePath, err)
		}

		return nil
	})
}

func mdPathToHTMLPath(mdPath string) string {
	return strings.TrimSuffix(mdPath, ".md") + ".html"
}
