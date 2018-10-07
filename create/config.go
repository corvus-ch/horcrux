package create

import (
	"io"

	"github.com/corvus-ch/horcrux/format"
)

type Config interface {
	Input() (io.Reader, error)
	Formats() ([]format.Format, error)
	Encrypt() bool
	Parts() int
	Threshold() int
}
