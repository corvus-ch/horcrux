package input

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/corvus-ch/horcrux/hash"
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

	cs, _ := hash.NewHashForPath(name)

	return &file{fi, stem, cs}
}

type file struct {
	fileInfo os.FileInfo
	stem     string
	hash     *hash.Hash
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

// Size returns the inputs size in bytes.
func (i *file) Size() int64 {
	if i.fileInfo == nil {
		return -1
	}

	return i.fileInfo.Size()
}

// Checksum returns the files checksum calculated for the given algorithm.
func (i *file) Checksum(alg string) (string, error) {
	return i.hash.Sum(alg)
}
