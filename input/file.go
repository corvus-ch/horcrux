package input

import (
	"os"
	"path/filepath"
	"strings"
)

func NewFileInput(name, stem string) Input {
	var fi os.FileInfo
	if len(name) > 0 {
		fi, _ = os.Stat(name)
	}

	if len(stem) == 0 && fi != nil {
		b := filepath.Base(fi.Name())
		stem = strings.TrimSuffix(b, filepath.Ext(b))
	}

	return &file{fi, stem}
}

type file struct {
	fileInfo os.FileInfo
	stem     string
}

func (i *file) Name() string {
	if i.fileInfo == nil {
		return ""
	}

	return i.fileInfo.Name()
}

func (i *file) Path() string {
	if i.fileInfo == nil {
		return ""
	}

	p, _ := filepath.Abs(i.fileInfo.Name())

	return p
}

func (i *file) Stem() string {
	return i.stem
}

func (i *file) Size() int64 {
	if i.fileInfo == nil {
		return -1
	}

	return i.fileInfo.Size()
}
