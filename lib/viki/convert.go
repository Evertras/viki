package viki

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	wikiLinks, err := buildWikiLinkMap(input)
	if err != nil {
		return fmt.Errorf("failed to build wiki link map: %w", err)
	}

	pathFilter, err := generatePathFilter(c.config, input)
	if err != nil {
		return fmt.Errorf("failed to generate path filter: %w", err)
	}

	dirTreeRoot, err := buildDirTree(input, pathFilter)
	if err != nil {
		return fmt.Errorf("failed to build root directory tree: %w", err)
	}

	sidebar, err := renderSidebar(dirTreeRoot)
	if err != nil {
		return fmt.Errorf("failed to render sidebar: %w", err)
	}

	err = walkFsRoot(input, pathFilter, func(inputFilePath string, info os.FileInfo) error {
		if filepath.Ext(inputFilePath) != ".md" {
			return nil
		}

		content, err := afero.ReadFile(input, inputFilePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", inputFilePath, err)
		}

		content = convertWikilinks(content, wikiLinks)
		htmlContent := mdToHtml(content)
		outputFile, err := renderPage(renderPageInput{
			Title:       strings.TrimSuffix(info.Name(), ".md"),
			BodyHtml:    htmlContent,
			SidebarHtml: sidebar,
		})
		if err != nil {
			return fmt.Errorf("failed to render page for %s: %w", inputFilePath, err)
		}

		outputFilePath := mdPathToHtmlPath(inputFilePath)

		err = output.MkdirAll(filepath.Dir(outputFilePath), 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", outputFilePath, err)
		}

		err = afero.WriteFile(output, outputFilePath, outputFile, 0644)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", outputFilePath, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to generate pages: %w", err)
	}

	// Add a root index.html file if it doesn't already exist
	exists, err := afero.Exists(output, "index.html")
	if err != nil {
		return fmt.Errorf("failed to check for existing index.html: %w", err)
	}

	if !exists {
		indexContent, err := renderIndexToc(dirTreeRoot)
		if err != nil {
			return fmt.Errorf("failed to render index content: %w", err)
		}
		indexPage, err := renderPage(renderPageInput{
			Title:       "Viki Home",
			BodyHtml:    indexContent,
			SidebarHtml: sidebar,
		})
		if err != nil {
			return fmt.Errorf("failed to render index page: %w", err)
		}
		err = afero.WriteFile(output, "index.html", indexPage, 0644)
		if err != nil {
			return fmt.Errorf("failed to write index.html: %w", err)
		}
	}

	cssContent, err := generateThemeCss(ThemeCatppuccinFrappe())

	if err != nil {
		return fmt.Errorf("failed to generate theme css: %w", err)
	}

	err = afero.WriteFile(output, "theme.css", []byte(cssContent), 0644)

	if err != nil {
		return fmt.Errorf("failed to write theme css: %w", err)
	}

	err = addStaticAssets(output)

	if err != nil {
		return fmt.Errorf("failed to add static assets: %w", err)
	}

	return nil
}
