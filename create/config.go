package create

import (
	"io"

	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/horcrux/input"
)

type Config interface {
	Input() (io.Reader, error)
	InputInfo() input.Input
	Formats() ([]format.Format, error)
	Encrypt() bool
	Parts() int
	Threshold() int
}
