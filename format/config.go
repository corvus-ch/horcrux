package format

import (
	"io"

	"github.com/corvus-ch/horcrux/input"
)

// Config define the configuration input required for creating an output format.
type Config interface {
	Input() (io.Reader, error)
	InputInfo() input.Input
	Formats() ([]Format, error)
}
