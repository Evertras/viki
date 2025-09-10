package viki

import (
	"fmt"
	"os"
	"path/filepath"
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
func (c *Converter) Convert(input afero.Fs, inputRootPath string, output afero.Fs, outputRootPath string) error {
	wikiLinks, err := c.buildWikiLinkMap(input, inputRootPath)
	if err != nil {
		return fmt.Errorf("failed to build wiki link map: %w", err)
	}

	sidebar, err := renderSidebar(input, inputRootPath)

	if err != nil {
		return fmt.Errorf("failed to render sidebar: %w", err)
	}

	err = afero.Walk(input, inputRootPath, func(inputFilePath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %s: %w", inputFilePath, err)
		}
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(inputFilePath) != ".md" {
			return nil
		}

		content, err := afero.ReadFile(input, inputFilePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", inputFilePath, err)
		}

		content = convertWikilinks(content, wikiLinks)
		content = mdToHtml(content)
		content = renderPage("Viki", string(content), sidebar)

		outputFilePath := mdPathToHTMLPath(inputFilePath)
		outputFilePath = strings.TrimPrefix(outputFilePath, inputRootPath)
		outputFilePath = filepath.Join(outputRootPath, outputFilePath)

		err = afero.WriteFile(output, outputFilePath, content, 0644)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", outputFilePath, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to generate pages: %w", err)
	}

	cssContent, err := c.generateThemeCss(ThemeCatpuccin())

	if err != nil {
		return fmt.Errorf("failed to generate theme css: %w", err)
	}

	err = afero.WriteFile(output, filepath.Join(outputRootPath, "theme.css"), []byte(cssContent), 0644)

	if err != nil {
		return fmt.Errorf("failed to write theme css: %w", err)
	}

	return nil
}

func mdPathToHTMLPath(mdPath string) string {
	mdPath = filepath.ToSlash(mdPath)
	return strings.TrimSuffix(mdPath, ".md") + ".html"
}
