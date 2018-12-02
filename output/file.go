package output

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/corvus-ch/horcrux/hash"
)

// File represents an output file.
type File interface {
	// Name returns the base name of the output file.
	Name() string

	// Path returns the absolute path of the output file.
	Path() string

	// Size returns the size of the output file in bytes.
	Size() int64

	// Checksum returns the outputs checksum calculated with the given algorithm.
	Checksum(alg string) (string, error)

	// Meta returns the mata data attached with the file for the given key.
	Meta(key string) (string, error)
}

type file struct {
	path string
	size int64
	hash *hash.Hash
	meta map[string]interface{}
}

// NewFile creates a new instance of File for the given path.
func NewFile(path string, meta map[string]interface{}) File {
	return &file{
		path: path,
		meta: meta,
		size: -1,
	}
}

// Name returns the outputs file name.
func (f *file) Name() string {
	return filepath.Base(f.path)
}

// Path returns the outputs file path.
func (f *file) Path() string {
	return f.path
}

// Size returns the outputs size. If the size can not be determined, the value will be -1.
func (f *file) Size() int64 {
	if f.size < 0 {
		if fi, err := os.Stat(f.path); err == nil {
			f.size = fi.Size()
		}
	}

	return f.size
}

// Checksum returns the outputs checksum calculated with the given algorithm.
// If the checksum can not be calculated, the returned value will be 'undefined'.
func (f *file) Checksum(alg string) (string, error) {
	if f.hash == nil {
		cs, err := hash.NewHashForPath(f.path)
		if err != nil {
			return "undefined", err
		}
		f.hash = cs
	}
	return f.hash.Sum(alg)
}

// Meta returns a piece of meta data attached to the output and identified by the given key.
func (f *file) Meta(key string) (string, error) {
	if f.meta == nil {
		return "", metadataKeyUnknown(key)
	}

	val, ok := f.meta[key]

	if !ok {
		return "", metadataKeyUnknown(key)
	}

	return fmt.Sprint(val), nil
}

func metadataKeyUnknown(key string) error {
	return fmt.Errorf("no meta data is known with key %s", key)
}