package input

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// NewFileInput creates an instance of Input representing a file input.
func NewFileInput(name, stem string) Input {
	var fi os.FileInfo
	if len(name) > 0 {
		fi, _ = os.Stat(name)
	}

	if len(stem) == 0 && fi != nil {
		b := filepath.Base(fi.Name())
		stem = strings.TrimSuffix(b, filepath.Ext(b))
	}

	return &file{fi, stem, createChecksums(name)}
}

type file struct {
	fileInfo   os.FileInfo
	stem       string
	checksumms *Hash
}

// Name returns the files name without its path.
func (i *file) Name() string {
	if i.fileInfo == nil {
		return ""
	}

	return i.fileInfo.Name()
}

// Path returns the files path including the file name.
func (i *file) Path() string {
	if i.fileInfo == nil {
		return ""
	}

	p, _ := filepath.Abs(i.fileInfo.Name())

	return p
}

// Path returns the inputs stem.
// By default, this is files name without the extension or whatever was
// set by the output option.
func (i *file) Stem() string {
	return i.stem
}

func (i *file) Size() int64 {
	if i.fileInfo == nil {
		return -1
	}

	return i.fileInfo.Size()
}

// Checksums returns a set containing the files checksums calculated
// with several algorithms.
func (i *file) Checksums() *Hash {
	return i.checksumms
}

func createChecksums(name string) *Hash {
	cs := NewHash()
	defer cs.Close()
	f, _ := os.Open(name)
	defer f.Close()
	if f != nil {
		io.Copy(cs, f)
	}
	return cs
}
