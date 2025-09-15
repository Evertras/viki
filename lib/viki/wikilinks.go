package viki

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/afero"
)

func buildWikiLinkMap(input afero.Fs) (map[string]string, error) {
	wikiLinkMap := make(map[string]string)

	err := afero.Walk(input, "", func(inputFilePath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %s: %w", inputFilePath, err)
		}
		if info.IsDir() {
			return nil
		}
		if path.Ext(inputFilePath) != ".md" {
			return nil
		}

		name := strings.TrimSuffix(info.Name(), ".md")

		// If we're running in Windows, convert backslashes to forward slashes
		if os.PathSeparator == '\\' {
			inputFilePath = strings.ReplaceAll(inputFilePath, "\\", "/")
		}

		relativePath := strings.TrimPrefix(inputFilePath, "/")
		relativePath = strings.TrimSuffix(relativePath, ".md") + ".html"

		// Make sure to start with a slash
		if !strings.HasPrefix(relativePath, "/") {
			relativePath = "/" + relativePath
		}
		wikiLinkMap[name] = relativePath

		return nil
	})

	return wikiLinkMap, err
}

func convertWikilinks(content []byte, wikiLinkMap map[string]string) []byte {
	// Find any text that matches the pattern [[Some Note]]
	// and replace it with [Some Note](Some%20Note.html) or whatever the appropriate link is
	// based on the wikiLinkMap.
	// TODO: This is inefficient for large link maps, but start here
	for wikiLink, target := range wikiLinkMap {
		linkText := fmt.Sprintf("[[%s]]", wikiLink)
		content = bytes.ReplaceAll(content, []byte(linkText),
			fmt.Appendf(nil, `<span class="wikilink">[%s](%s)</span>`, wikiLink, target))
	}

	return content
}
