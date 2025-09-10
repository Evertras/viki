package viki

import (
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestBuildWikiLinkMap(t *testing.T) {
	converter := NewConverter(ConverterOptions{})
	assert.NotNil(t, converter)
	inputFs := afero.NewMemMapFs()

	filePaths := []string{
		"/wiki/note1.md",
		"/wiki/subdir/note2.md",
		"/wiki/note3.txt",   // Non-md file, should be ignored
		"/another/note4.md", // Outside the wiki directory, should be ignored
	}

	expectedMap := map[string]string{
		"note1": "/" + path.Join("note1.html"),
		"note2": "/" + path.Join("subdir", "note2.html"),
	}

	for _, filePath := range filePaths {
		afero.WriteFile(inputFs, filePath, []byte("test content"), 0644)
	}

	wikiLinkMap, err := converter.buildWikiLinkMap(inputFs, "/wiki")
	assert.NoError(t, err)
	assert.Equal(t, expectedMap, wikiLinkMap)
}

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
