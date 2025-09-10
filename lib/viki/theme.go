package viki

import (
	"bytes"
	_ "embed"

	"html/template"

	catppuccin "github.com/catppuccin/go"
)

var (
	//go:embed templates/base/theme.css.tpl
	themeCssTemplateRaw []byte

	themeCssTemplate *template.Template
)

type ThemeData struct {
	BgColor               string
	FgColor               string
	LinkColor             string
	LinkHoverColor        string
	StrongColor           string
	HeaderColor           string
	SidebarBgColor        string
	SidebarFgColor        string
	CodeBgColor           string
	CodeFgColor           string
	BlockquoteBgColor     string
	BlockquoteBorderColor string
}

func init() {
	themeCssTemplate = template.Must(template.New("theme").Parse(string(themeCssTemplateRaw)))
}

func (c *Converter) generateThemeCss(data ThemeData) ([]byte, error) {
	var buf bytes.Buffer
	if err := themeCssTemplate.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ThemeCatpuccin() ThemeData {
	flavor := catppuccin.Frappe
	return ThemeData{
		BgColor:               flavor.Base().Hex,
		FgColor:               flavor.Text().Hex,
		LinkColor:             flavor.Blue().Hex,
		LinkHoverColor:        flavor.Lavender().Hex,
		StrongColor:           flavor.Teal().Hex,
		HeaderColor:           flavor.Green().Hex,
		SidebarBgColor:        flavor.Mantle().Hex,
		SidebarFgColor:        flavor.Text().Hex,
		CodeBgColor:           flavor.Base().Hex,
		CodeFgColor:           flavor.Text().Hex,
		BlockquoteBgColor:     flavor.Base().Hex,
		BlockquoteBorderColor: flavor.Overlay1().Hex,
	}
}
