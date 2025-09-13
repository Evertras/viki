package viki

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/spf13/afero"
)

type ConverterOptions struct {
	ExcludePatterns []string
	IncludePatterns []string
}

type Converter struct {
	config ConverterOptions
}

func NewConverter(options ConverterOptions) *Converter {
	return &Converter{
		config: options,
	}
}

// Convert takes in a filesystem converts each .md file to .html, and then writes to
// the output filesystem. It always reads and writes from the root, so always use
// afero.NewBasePathFs to scope it to a subdirectory.
func (c *Converter) Convert(input afero.Fs, output afero.Fs) error {
	// Enforce basepathfs being used
	_, ok := input.(*afero.BasePathFs)
	if !ok {
		return fmt.Errorf("input filesystem must be a BasePathFs for safety")
	}
	_, ok = output.(*afero.BasePathFs)
	if !ok {
		return fmt.Errorf("output filesystem must be a BasePathFs for safety")
	}

	// Be extra safe and make input read-only
	input = afero.NewReadOnlyFs(input)

	wikiLinks, err := c.buildWikiLinkMap(input)
	if err != nil {
		return fmt.Errorf("failed to build wiki link map: %w", err)
	}

	ignoreChecker, err := c.generateIgnoreChecker(input)

	if err != nil {
		return fmt.Errorf("failed to generate ignore checker: %w", err)
	}

	if len(c.config.IncludePatterns) == 0 {
		c.config.IncludePatterns = []string{"**"}
	}

	includeChecker := ignore.CompileIgnoreLines(c.config.IncludePatterns...)

	sidebar, err := renderSidebar(input, ignoreChecker, includeChecker)

	if err != nil {
		return fmt.Errorf("failed to render sidebar: %w", err)
	}

	err = afero.Walk(input, "", func(inputFilePath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %s: %w", inputFilePath, err)
		}

		if info.IsDir() ||
			ignoreChecker.MatchesPath(inputFilePath) ||
			!includeChecker.MatchesPath(inputFilePath) {
			return nil
		}

		if ignoreChecker.MatchesPath(inputFilePath) {
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

		err = output.MkdirAll(filepath.Dir(outputFilePath), 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", outputFilePath, err)
		}

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

	err = afero.WriteFile(output, "theme.css", []byte(cssContent), 0644)

	if err != nil {
		return fmt.Errorf("failed to write theme css: %w", err)
	}

	err = c.addStaticAssets(output)

	if err != nil {
		return fmt.Errorf("failed to add static assets: %w", err)
	}

	return nil
}

func mdPathToHTMLPath(mdPath string) string {
	mdPath = filepath.ToSlash(mdPath)
	return strings.TrimSuffix(mdPath, ".md") + ".html"
}
