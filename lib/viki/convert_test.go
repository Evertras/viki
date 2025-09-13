package viki

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestConverterDoesNothingFromEmptyFs(t *testing.T) {
	converter := NewConverter(ConverterOptions{})
	assert.NotNil(t, converter)

	inputFs := afero.NewBasePathFs(afero.NewMemMapFs(), "/")
	outputFs := afero.NewBasePathFs(afero.NewMemMapFs(), "/")

	err := converter.Convert(inputFs, outputFs)
	assert.NoError(t, err)

	// Verify that no files were created in the output filesystem
	afero.Walk(outputFs, "", func(outputFilePath string, info os.FileInfo, err error) error {
		assert.NoError(t, err, "verification walk function should not error")

		if err != nil {
			return err
		}

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

	inputFs := afero.NewBasePathFs(afero.NewMemMapFs(), "/")
	outputFs := afero.NewBasePathFs(afero.NewMemMapFs(), "/")

	// Create a non-md file in the input filesystem
	afero.WriteFile(inputFs, "test.txt", []byte("test"), 0644)

	err := converter.Convert(inputFs, outputFs)
	assert.NoError(t, err)

	// Verify that no files were created in the output filesystem
	afero.Walk(outputFs, "", func(outputFilePath string, info os.FileInfo, err error) error {
		assert.NoError(t, err)

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		assert.False(t, info.IsDir(), "Expected no files to be created")
		return nil
	})
}

func TestConverterRespectsGitIgnore(t *testing.T) {
	converter := NewConverter(ConverterOptions{
		ExcludePatterns: []string{"/private"},
	})
	assert.NotNil(t, converter)

	inputFs := afero.NewBasePathFs(afero.NewMemMapFs(), "/")
	outputFs := afero.NewBasePathFs(afero.NewMemMapFs(), "/")

	writeFile := func(path string, content string) {
		err := afero.WriteFile(inputFs, path, []byte(content), 0644)
		assert.NoError(t, err)
	}

	// Create a non-md file in the input filesystem
	writeFile("test.txt", "test")
	writeFile("node_modules/thing/README.md", "# Hello this is a thing from node_modules")
	writeFile("private/private-thing.md", "# Hello this is a private thing")
	writeFile("public/public.md", "# Hello this is a public thing")
	writeFile(".gitignore", "node_modules")

	err := converter.Convert(inputFs, outputFs)
	assert.NoError(t, err)

	// Verify that no files were created in the output filesystem
	err = afero.Walk(outputFs, "", func(outputFilePath string, info os.FileInfo, err error) error {
		assert.NoError(t, err)

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(outputFilePath, ".html") {
			return nil
		}

		assert.Equal(t, "public.html", info.Name(), "Expected only public.md to be converted")

		return nil
	})

	assert.NoError(t, err, "Walk should not error")
}

func TestConverterIncludesOnlySpecifiedPatterns(t *testing.T) {
	converter := NewConverter(ConverterOptions{
		IncludePatterns: []string{"/included"},
		ExcludePatterns: []string{"excluded"},
	})
	assert.NotNil(t, converter)
	inputFs := afero.NewBasePathFs(afero.NewMemMapFs(), "/")
	outputFs := afero.NewBasePathFs(afero.NewMemMapFs(), "/")
	writeFile := func(path string, content string) {
		err := afero.WriteFile(inputFs, path, []byte(content), 0644)
		assert.NoError(t, err, "Failed to write file %s in test setup", path)
	}

	// Create a non-md file in the input filesystem
	writeFile("included/included-file.md", "# This file should be included")
	writeFile("excluded/excluded-file.md", "# This file should be excluded")
	writeFile("included/also-included.md", "# This file should also be included")
	writeFile("not-included/not-included.md", "# This file should not be included")
	writeFile("included/excluded/nope.md", "# This file should be excluded as well")

	err := converter.Convert(inputFs, outputFs)
	assert.NoError(t, err, "Conversion should not error")

	seenFiles := []string{}
	err = afero.Walk(outputFs, "", func(outputFilePath string, info os.FileInfo, err error) error {
		assert.NoError(t, err, "Walk function should not error")

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), ".html") {
			seenFiles = append(seenFiles, info.Name())
		}
		return nil
	})

	assert.NoError(t, err, "Walk should not error")
	assert.ElementsMatch(t, []string{"included-file.html", "also-included.html"}, seenFiles, "Expected only included files to be created")
}

func TestConverterCreatesFilesWithSameNameButHtmlExtension(t *testing.T) {
	converter := NewConverter(ConverterOptions{})
	assert.NotNil(t, converter)
	inputFs := afero.NewBasePathFs(afero.NewMemMapFs(), "/")
	outputFs := afero.NewBasePathFs(afero.NewMemMapFs(), "/")

	// Create a .md file in the input filesystem
	afero.WriteFile(inputFs, "Test.md", []byte("# Test"), 0644)
	err := converter.Convert(inputFs, afero.NewBasePathFs(outputFs, "/site"))
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

func TestConverterAddsStaticAssets(t *testing.T) {
	converter := NewConverter(ConverterOptions{})
	assert.NotNil(t, converter)
	outputFs := afero.NewMemMapFs()
	err := converter.addStaticAssets(afero.NewBasePathFs(outputFs, "/site"))
	assert.NoError(t, err, "Adding static assets should not error")
	// Verify that static assets were added
	for path := range staticAssetFileMap {
		fullPath := "/site/" + path
		exists, err := afero.Exists(outputFs, fullPath)
		assert.NoError(t, err, "Existence check should not error")
		assert.True(t, exists, "Expected %s to exist in the output filesystem", fullPath)
	}
}

func TestConverterWorksForHttpServe(t *testing.T) {
	converter := NewConverter(ConverterOptions{})
	assert.NotNil(t, converter)
	inputFs := afero.NewBasePathFs(afero.NewMemMapFs(), "/")
	outputFs := afero.NewBasePathFs(afero.NewMemMapFs(), "/")

	// Create a .md file in the input filesystem
	afero.WriteFile(inputFs, "about.md", []byte("# About Page"), 0644)

	err := converter.Convert(inputFs, outputFs)
	assert.NoError(t, err, "Conversion should not error")

	// Verify that the corresponding .html file was created in the output filesystem
	exists, err := afero.Exists(outputFs, "about.html")
	assert.NoError(t, err, "Existence check should not error")
	assert.True(t, exists, "Expected about.html to exist in the output filesystem")

	httpFs := afero.NewHttpFs(outputFs)
	handler := http.FileServer(httpFs.Dir(""))

	req := httptest.NewRequest("GET", "http://example.com/about.html", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status OK for about.html")
}
