package viki

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/afero"
)

func renderSidebar(fs afero.Fs, pathFilter pathFilter) (template.HTML, error) {
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

	err := walkFsRoot(fs, pathFilter, func(filePath string, info os.FileInfo) error {
		if info.IsDir() {
			nodes[filePath] = &node{
				Name:     info.Name(),
				URL:      filepathToEscapedHttpPath(filePath),
				IsDir:    true,
				Children: []*node{},
			}

			parentDir := filepath.Dir(filePath)

			parent := nodes[parentDir]
			if parent == nil {
				parent = &node{
					Name:     parentDir,
					URL:      filepathToEscapedHttpPath(parentDir),
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
				URL:   mdPathToHtmlPath(filePath),
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

	// Put directories above individual pages, and sort alphabetically
	var sortChildren func(n *node)
	sortChildren = func(n *node) {
		slices.SortFunc(n.Children, func(a, b *node) int {
			if a.IsDir && !b.IsDir {
				return -1
			}
			if !a.IsDir && b.IsDir {
				return 1
			}
			return strings.Compare(a.Name, b.Name)
		})

		for _, child := range n.Children {
			sortChildren(child)
		}
	}

	sortChildren(rootNode)

	err = template_base_sidebar_gohtml.Execute(&out, map[string]any{
		"Nodes": rootNode.Children,
	})

	if err != nil {
		return "", fmt.Errorf("failed to render sidebar template: %w", err)
	}

	return template.HTML(out.String()), nil
}
