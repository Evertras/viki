package viki

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

func (c *Converter) addStaticAssets(fs afero.Fs, outputPath string) error {
	staticDirPath := filepath.Join(outputPath, "_viki_static")
	staticDir := afero.NewBasePathFs(fs, staticDirPath)

	err := staticDir.MkdirAll(staticDirPath, 0755)

	if err != nil {
		return err
	}

	for path, data := range staticAssetFileMap {
		err := afero.WriteFile(staticDir, path, data, 0644)

		if err != nil {
			return err
		}
	}

	// Special case for favicon
	err = afero.WriteFile(fs, filepath.Join(outputPath, "favicon.ico"), static_favicon_favicon_ico, 0644)

	if err != nil {
		return fmt.Errorf("failed to write favicon: %w", err)
	}

	return nil
}
