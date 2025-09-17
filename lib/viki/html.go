package viki

import (
	"bytes"
	"fmt"
	"html/template"
	"io"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func mdToHtml(mdContent []byte) template.HTML {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.Footnotes
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(mdContent)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: mdToHtmlRenderHook,
	}
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

func renderTocFromTemplate(rootNode *dirTreeNode, tpl *template.Template) (template.HTML, error) {
	var out bytes.Buffer

	var nodes []*dirTreeNode
	if rootNode != nil {
		nodes = rootNode.Children
	}

	err := tpl.Execute(&out, map[string]any{
		"Nodes": nodes,
	})

	if err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return template.HTML(out.String()), nil
}

func renderSidebar(rootNode *dirTreeNode) (template.HTML, error) {
	return renderTocFromTemplate(rootNode, template_base_sidebar_gohtml)
}

func renderIndexToc(rootNode *dirTreeNode) (template.HTML, error) {
	return renderTocFromTemplate(rootNode, template_base_index_gohtml)
}

// References:
//   - https://blog.kowalczyk.info/article/cxn3/advanced-markdown-processing-in-go.html
//   - https://github.com/alecthomas/chroma/blob/master/quick/quick.go
func renderCode(w io.Writer, codeBlock *ast.CodeBlock, entering bool) {
	lang := string(codeBlock.Info)
	err := quick.Highlight(w, string(codeBlock.Literal), lang, "html", "catppuccin-frappe")

	if err != nil {
		// Fallback: just render the code block as-is
		if entering {
			fmt.Fprintf(w, "<pre><code>")
			template.HTMLEscape(w, codeBlock.Literal)
			fmt.Fprintf(w, "</code></pre>")
		}
	}
}

func mdToHtmlRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if code, ok := node.(*ast.CodeBlock); ok {
		renderCode(w, code, entering)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}
