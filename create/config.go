package create

import (
	"io"

	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/horcrux/input"
)

type Config interface {
	Reader() (io.Reader, error)
	InputInfo() input.Input
	Formats() ([]format.Format, error)
	Encrypted() bool
	Parts() int
	Threshold() int
}
