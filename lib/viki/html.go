package viki

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func mdToHtml(mdContent []byte) template.HTML {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.Footnotes
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(mdContent)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	htmlContent := markdown.Render(doc, renderer)

	return template.HTML(htmlContent)
}

type renderPageInput struct {
	Title       string
	BodyHtml    template.HTML
	SidebarHtml template.HTML
}

func renderPage(input renderPageInput) ([]byte, error) {
	var out bytes.Buffer
	err := template_base_page_gohtml.Execute(&out, input)
	if err != nil {
		return nil, fmt.Errorf("failed to render page template: %w", err)
	}
	return out.Bytes(), nil
}

func renderSidebar(rootNode *dirTreeNode, pathFilter pathFilter) (template.HTML, error) {
	var out bytes.Buffer

	if len(rootNode.Children) == 0 {
		return "No content", nil
	}

	err := template_base_sidebar_gohtml.Execute(&out, map[string]any{
		"Nodes": rootNode.Children,
	})

	if err != nil {
		return "", fmt.Errorf("failed to render sidebar template: %w", err)
	}

	return template.HTML(out.String()), nil
}

func renderIndex(rootNode *dirTreeNode) (template.HTML, error) {
	var out bytes.Buffer

	if len(rootNode.Children) == 0 {
		return "No content", nil
	}

	err := template_base_index_gohtml.Execute(&out, map[string]any{
		"Nodes": rootNode.Children,
	})

	if err != nil {
		return "", fmt.Errorf("failed to render index template: %w", err)
	}
	return template.HTML(out.String()), nil
}
