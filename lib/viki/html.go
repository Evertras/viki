package viki

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	ignore "github.com/sabhiram/go-gitignore"
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

func renderSidebar(fs afero.Fs, ignoreChecker *ignore.GitIgnore, includeChecker *ignore.GitIgnore) (string, error) {
	type node struct {
		Name     string
		URL      string
		IsDir    bool
		Children []*node
	}

	var out bytes.Buffer

	nodes := make(map[string]*node)

	rootNode := &node{
		Name:     "Root",
		URL:      "/",
		IsDir:    true,
		Children: []*node{},
	}

	nodes["."] = rootNode

	err := afero.Walk(fs, "", func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filePath == "." ||
			ignoreChecker.MatchesPath(filePath) ||
			!includeChecker.MatchesPath(filePath) {
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
				Name:  strings.TrimSuffix(info.Name(), ".md"),
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
	var hasAnyChildren func(n *node) bool
	hasAnyChildren = func(n *node) bool {
		if !n.IsDir {
			return true
		}
		return slices.ContainsFunc(n.Children, hasAnyChildren)
	}

	var pruneChildren func(n *node)
	pruneChildren = func(n *node) {
		if !n.IsDir {
			return
		}
		for i := len(n.Children) - 1; i >= 0; i-- {
			if !hasAnyChildren(n.Children[i]) {
				n.Children = append(n.Children[:i], n.Children[i+1:]...)
			} else {
				pruneChildren(n.Children[i])
			}
		}
	}

	pruneChildren(rootNode)

	if len(rootNode.Children) == 0 {
		return "No content", nil
	}

	err = sidebarTemplate.Execute(&out, map[string]any{
		"Nodes": rootNode.Children,
	})

	if err != nil {
		return "", fmt.Errorf("failed to render sidebar template: %w", err)
	}

	return out.String(), nil
}

func renderPage(title, body, sidebar string) []byte {
	var out bytes.Buffer
	pageTemplate.Execute(&out, map[string]any{
		"Title":       title,
		"BodyHtml":    template.HTML(body),
		"SidebarHtml": template.HTML(sidebar),
	})
	return out.Bytes()
}
