package assert

import "github.com/corvus-ch/horcrux/format"

// FormatFactory describes a func used for instantiating a Format during assertions.
type FormatFactory func(string) format.Format
