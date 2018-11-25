package text

import (
	"fmt"

	"github.com/corvus-ch/horcrux/format/qr"
	"github.com/corvus-ch/horcrux/input"
	"github.com/corvus-ch/horcrux/output"
)

type Data struct {
	Input input.Input
	Lines chan Line
	files map[string][]output.File
}

func NewData(in input.Input, lines chan Line) *Data {
	d := &Data{
		Input: in,
		Lines: lines,
		files: make(map[string][]output.File, 0),
	}

	d.AppendFile(qr.Name, "test.019.1.png")
	d.AppendFile(qr.Name, "test.161.1.png")
	d.AppendFile(qr.Name, "test.182.1.png")

	return d
}

func (d *Data) OutputFiles(format string) (files []output.File, err error) {
	files, ok := d.files[format]
	if !ok {
		err = fmt.Errorf("unknown format %s", format)
	}
	return
}

func (d *Data) AppendFile(format, path string) {
	files, ok := d.files[format]
	if ok {
		d.files[format] = append(files, output.New(path))
	} else {
		d.files[format] = []output.File{output.New(path)}
	}
}
