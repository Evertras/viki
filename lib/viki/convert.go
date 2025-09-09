package viki

import "io/fs"

type ConverterOptions struct {
}

type Converter struct {
}

func NewConverter(options ConverterOptions) *Converter {
	return &Converter{}
}

// Convert takes in a filesystem (which could be an OS filesystem or an in-memory one)
// and an output path, and converts the contents of the filesystem into an in-memory
// file system ready to be written out. Does not remove anything from the output fs,
// just adds/overwrites files, so give it a fresh one.
func (c *Converter) Convert(fs fs.FS, output fs.FS) error {
	return nil
}
