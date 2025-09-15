package viki

import (
	"bytes"

	catppuccin "github.com/catppuccin/go"
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
	ListBulletColor       string
}

func generateThemeCss(data ThemeData) ([]byte, error) {
	var buf bytes.Buffer
	if err := template_base_theme_css_tpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ThemeCatppuccinFrappe() ThemeData {
	flavor := catppuccin.Frappe

	return ThemeData{
		BgColor:               flavor.Base().Hex,
		FgColor:               flavor.Text().Hex,
		LinkColor:             flavor.Blue().Hex,
		LinkHoverColor:        flavor.Lavender().Hex,
		StrongColor:           flavor.Peach().Hex,
		HeaderColor:           flavor.Green().Hex,
		SidebarBgColor:        flavor.Mantle().Hex,
		SidebarFgColor:        flavor.Text().Hex,
		CodeBgColor:           flavor.Base().Hex,
		CodeFgColor:           flavor.Text().Hex,
		BlockquoteBgColor:     flavor.Base().Hex,
		BlockquoteBorderColor: flavor.Peach().Hex,
		ListBulletColor:       flavor.Green().Hex,
	}
}
