package restore

import (
	"io"

	"github.com/corvus-ch/horcrux/format"
)

type Config interface {
	Format() (format.Format, error)
	Decrypt() bool
	Output() (io.Writer, error)
	FileNames() []string
}
