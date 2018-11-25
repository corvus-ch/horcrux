package text

import (
	"fmt"

	"github.com/corvus-ch/horcrux/format/qr"
	"github.com/corvus-ch/horcrux/input"
)

type Data struct {
	Input input.Input
	Lines chan Line
	files map[string][]OutputFile
}

type OutputFile struct {
	Path string
}

func NewData(in input.Input, lines chan Line) *Data {
	d := &Data{
		Input: in,
		Lines: lines,
		files: make(map[string][]OutputFile, 0),
	}

	d.AppendFile(qr.Name, "test.019.1.png")
	d.AppendFile(qr.Name, "test.161.1.png")
	d.AppendFile(qr.Name, "test.182.1.png")

	return d
}

func (d *Data) OutputFiles(format string) (files []OutputFile, err error) {
	files, ok := d.files[format]
	if !ok {
		err = fmt.Errorf("unknown format %s", format)
	}
	return
}

func (d *Data) AppendFile(format, path string) {
	files, ok := d.files[format]
	if ok {
		d.files[format] = append(files, OutputFile{path})
	} else {
		d.files[format] = []OutputFile{{path}}
	}
}
