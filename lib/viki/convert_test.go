package viki

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestConverterDoesNothingFromEmptyFs(t *testing.T) {
	converter := NewConverter(ConverterOptions{})
	assert.NotNil(t, converter)

	inputFs := afero.NewMemMapFs()
	outputFs := afero.NewMemMapFs()

	err := converter.Convert(inputFs, "/", outputFs, "/")
	assert.NoError(t, err)

	// Verify that no files were created in the output filesystem
	afero.Walk(outputFs, "/", func(outputFilePath string, info os.FileInfo, err error) error {
		assert.NoError(t, err)

		if info.IsDir() {
			return nil
		}

		assert.False(t, info.IsDir(), "Expected no files to be created")
		return nil
	})
}

func TestConverterDoesNothingToNonMdFiles(t *testing.T) {
	converter := NewConverter(ConverterOptions{})
	assert.NotNil(t, converter)

	inputFs := afero.NewMemMapFs()
	outputFs := afero.NewMemMapFs()

	// Create a non-md file in the input filesystem
	afero.WriteFile(inputFs, "/test.txt", []byte("test"), 0644)

	err := converter.Convert(inputFs, "/", outputFs, "/")
	assert.NoError(t, err)

	// Verify that no files were created in the output filesystem
	afero.Walk(outputFs, "/", func(outputFilePath string, info os.FileInfo, err error) error {
		assert.NoError(t, err)

		if info.IsDir() {
			return nil
		}

		assert.False(t, info.IsDir(), "Expected no files to be created")
		return nil
	})
}

func TestConverterCreatesFilesWithSameNameButHtmlExtension(t *testing.T) {
	converter := NewConverter(ConverterOptions{})
	assert.NotNil(t, converter)
	inputFs := afero.NewMemMapFs()
	outputFs := afero.NewMemMapFs()

	// Create a .md file in the input filesystem
	afero.WriteFile(inputFs, "/Test.md", []byte("# Test"), 0644)
	err := converter.Convert(inputFs, "/", outputFs, "/site")
	assert.NoError(t, err, "Conversion should not error")

	// Verify that the corresponding .html file was created in the output filesystem
	exists, err := afero.Exists(outputFs, "/site/Test.html")
	assert.NoError(t, err, "Existence check should not error")
	assert.True(t, exists, "Expected /site/Test.html to exist in the output filesystem")
}

func TestMDPathToHTMLPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"test.md", "test.html"},
		{"another/test.md", "another/test.html"},
		{"no_extension", "no_extension.html"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := mdPathToHTMLPath(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
