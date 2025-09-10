package viki

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/spf13/afero"
)

var (
	//go:embed templates/base/sidebar.html.tpl
	sidebarTemplateRaw []byte

	//go:embed templates/base/page.html.tpl
	pageTemplateRaw []byte

	sidebarTemplate *template.Template
	pageTemplate    *template.Template
)

func init() {
	sidebarTemplate = template.Must(template.New("base").Parse(string(sidebarTemplateRaw)))
	pageTemplate = template.Must(template.New("base").Parse(string(pageTemplateRaw)))
}

func mdToHtml(mdContent []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.Footnotes
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(mdContent)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	htmlContent := markdown.Render(doc, renderer)

	return htmlContent
}

func renderSidebar(fs afero.Fs, basePath string) (string, error) {
	type node struct {
		Name     string
		URL      string
		IsDir    bool
		Children []*node
	}

	var out bytes.Buffer

	nodes := make(map[string]*node)

	nodes["."] = &node{
		Name:     "Root",
		URL:      "/",
		IsDir:    true,
		Children: []*node{},
	}

	err := afero.Walk(fs, basePath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filePath[0] == '.' {
			return nil
		}

		if info.IsDir() {
			nodes[filePath] = &node{
				Name:     info.Name(),
				URL:      filePath,
				IsDir:    true,
				Children: []*node{},
			}

			parentDir := filepath.Dir(filePath)

			parent := nodes[parentDir]
			if parent == nil {
				parent = &node{
					Name:     parentDir,
					URL:      parentDir,
					IsDir:    true,
					Children: []*node{},
				}
			}
			parent.Children = append(parent.Children, nodes[filePath])
			nodes[parentDir] = parent

			return nil
		}

		// Add .md files to children of their parent directory
		if filepath.Ext(filePath) == ".md" {
			parentDir := filepath.Dir(filePath)

			parent := nodes[parentDir]
			parent.Children = append(parent.Children, &node{
				Name:  info.Name(),
				URL:   mdPathToHTMLPath(filePath),
				IsDir: false,
			})
			nodes[parentDir] = parent
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to walk filesystem for sidebar: %w", err)
	}

	// Prune any empty directories
	// TODO: Do this better so we don't leave empty parents
	for filePath, n := range nodes {
		if n.IsDir && len(n.Children) == 0 {
			delete(nodes, filePath)
		}
	}

	if nodes["."] == nil {
		return "No content", nil
	}

	err = sidebarTemplate.Execute(&out, map[string]any{
		"Nodes": nodes["."].Children,
	})

	if err != nil {
		return "", fmt.Errorf("failed to render sidebar template: %w", err)
	}

	return out.String(), nil
}

func renderPage(body, sidebar string) []byte {
	var out bytes.Buffer
	pageTemplate.Execute(&out, map[string]any{
		"BodyHtml":    template.HTML(body),
		"SidebarHtml": template.HTML(sidebar),
	})
	return out.Bytes()
}
