package text

import (
	"fmt"

	"github.com/corvus-ch/horcrux/input"
	"github.com/corvus-ch/horcrux/output"
)

// Data holds the values passed to the template.
type Data struct {
	Input input.Input
	Lines chan Line
	X     byte
	files output.Output
}

// NewData creates a new Daa instance.
func NewData(in input.Input, x byte, out output.Output) *Data {
	return &Data{
		Input: in,
		Lines: make(chan Line),
		X:     x,
		files: out,
	}
}

// OutputFiles returns the list of generates files for a given output format.
//
// Since the list might not be known at the time this function gets called, a channel is returned.
// Files will be written to that channel and the channel will be closed when no new files have to be expected.
func (d *Data) OutputFiles(format string) (files chan output.File, err error) {
	files, ok := d.files[format]
	if !ok {
		err = fmt.Errorf("unknown format %s", format)
	}
	return
}
