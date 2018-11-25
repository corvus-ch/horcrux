package text

import (
	"fmt"

	"github.com/corvus-ch/horcrux/input"
	"github.com/corvus-ch/horcrux/output"
)

type Data struct {
	Input input.Input
	Lines chan Line
	X     byte
	files output.Output
}

func NewData(in input.Input, x byte, out output.Output) *Data {
	d := &Data{
		Input: in,
		Lines: make(chan Line),
		X:     x,
		files: out,
	}

	return d
}

func (d *Data) OutputFiles(format string) (files chan output.File, err error) {
	files, ok := d.files[format]
	if !ok {
		err = fmt.Errorf("unknown format %s", format)
	}
	return
}
