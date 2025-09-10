package viki

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWikilinkConversion(t *testing.T) {
	wikiLinkMap := map[string]string{
		"Note One": "Note%20One.html",
		"Note Two": "Note%20Two.html",
	}
	originalContent := []byte("This is a link to [[Note One]] and another link to [[Note Two]]. [[This link]] doesn't exist")
	expectedContent := []byte("This is a link to [Note One](Note%20One.html) and another link to [Note Two](Note%20Two.html). [[This link]] doesn't exist")
	convertedContent := convertWikilinks(originalContent, wikiLinkMap)
	assert.Equal(t, string(expectedContent), string(convertedContent))
}
