package restore

import (
	"io"

	"github.com/corvus-ch/horcrux/format"
)

// Config holds the values controlling the restore behaviour.
type Config interface {
	Format() (format.Format, error)
	Decrypt() bool
	Output() (io.Writer, error)
	FileNames() []string
}
