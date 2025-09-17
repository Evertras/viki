package viki

import (
	"bytes"
	"fmt"
	"html/template"
	"io"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
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
func renderCode(w io.Writer, codeBlock *ast.CodeBlock) error {
	const style = "catppuccin-frappe"

	lang := string(codeBlock.Info)
	source := string(codeBlock.Literal)

	// Determine lexer.
	l := lexers.Get(lang)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	// Determine formatter.
	f := chromahtml.New(chromahtml.ClassPrefix("chroma-"), chromahtml.TabWidth(2), chromahtml.WithLineNumbers(true))

	// Determine style.
	s := styles.Get(style)
	if s == nil {
		s = styles.Fallback
	}

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return fmt.Errorf("failed to tokenize source: %w", err)
	}

	return f.Format(w, s, it)
}

func mdToHtmlRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if code, ok := node.(*ast.CodeBlock); ok {
		var out bytes.Buffer
		err := renderCode(&out, code)
		if err != nil {
			return ast.GoToNext, false
		}
		_, err = out.WriteTo(w)
		if err != nil {
			return ast.GoToNext, false
		}
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}
