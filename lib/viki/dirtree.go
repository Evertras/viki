package viki

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/afero"
)

type dirTreeNode struct {
	Name     string
	URL      string
	IsDir    bool
	Children []*dirTreeNode
}

func buildDirTree(fs afero.Fs, pathFilter pathFilter) (*dirTreeNode, error) {
	nodes := make(map[string]*dirTreeNode)

	rootNode := &dirTreeNode{
		Name:     "Root",
		URL:      "/",
		IsDir:    true,
		Children: []*dirTreeNode{},
	}

	nodes["."] = rootNode

	err := walkFsRoot(fs, pathFilter, func(filePath string, info os.FileInfo) error {
		if info.IsDir() {
			nodes[filePath] = &dirTreeNode{
				Name:     info.Name(),
				URL:      filepathToEscapedHttpPath(filePath),
				IsDir:    true,
				Children: []*dirTreeNode{},
			}

			parentDir := filepath.Dir(filePath)

			parent := nodes[parentDir]
			if parent == nil {
				parent = &dirTreeNode{
					Name:     parentDir,
					URL:      filepathToEscapedHttpPath(parentDir),
					IsDir:    true,
					Children: []*dirTreeNode{},
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
			parent.Children = append(parent.Children, &dirTreeNode{
				Name:  strings.TrimSuffix(info.Name(), ".md"),
				URL:   mdPathToHtmlPath(filePath),
				IsDir: false,
			})
			nodes[parentDir] = parent
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk filesystem for sidebar: %w", err)
	}

	// Prune any empty directories
	var hasAnyChildren func(n *dirTreeNode) bool
	hasAnyChildren = func(n *dirTreeNode) bool {
		if !n.IsDir {
			return true
		}
		return slices.ContainsFunc(n.Children, hasAnyChildren)
	}

	var pruneChildren func(n *dirTreeNode)
	pruneChildren = func(n *dirTreeNode) {
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

	// Put directories above individual pages, and sort alphabetically
	var sortChildren func(n *dirTreeNode)
	sortChildren = func(n *dirTreeNode) {
		slices.SortFunc(n.Children, func(a, b *dirTreeNode) int {
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

	return rootNode, nil
}
