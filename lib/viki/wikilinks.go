package viki

import (
	"bytes"
	"fmt"
)

func convertWikilinks(content []byte, wikiLinkMap map[string]string) []byte {
	// Find any text that matches the pattern [[Some Note]]
	// and replace it with [Some Note](Some%20Note.html) or whatever the appropriate link is
	// based on the wikiLinkMap.
	// TODO: This is inefficient for large link maps, but start here
	for wikiLink, target := range wikiLinkMap {
		linkText := fmt.Sprintf("[[%s]]", wikiLink)
		content = bytes.ReplaceAll(content, []byte(linkText), []byte(fmt.Sprintf("[%s](%s)", wikiLink, target)))
	}

	return content
}
