package viki

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/afero"
)

func (c *Converter) addStaticAssets(fs afero.Fs) error {
	for path, data := range staticAssetFileMap {
		dir := filepath.Dir(path)
		log.Println("Creating dir:", dir)
		err := fs.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to write dir %s for file %s: %w", dir, path, err)
		}

		err = afero.WriteFile(fs, path, data, 0644)

		if err != nil {
			return fmt.Errorf("failed to write static asset %s: %w", path, err)
		}
	}

	// Special case for favicon
	err := afero.WriteFile(fs, "favicon.ico", static_favicon_favicon_ico, 0644)

	if err != nil {
		return fmt.Errorf("failed to write favicon: %w", err)
	}

	return nil
}
