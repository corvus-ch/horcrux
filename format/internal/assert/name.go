package assert

import (
	"fmt"
	"testing"

	"github.com/corvus-ch/horcrux/format"
	"github.com/stretchr/testify/assert"
)

// NameTest holds the data required to assert the behaviour of github.com/corvus-ch/horcrux/format.Format.
type NameTest struct {
	X        byte
	Stem     string
	Expected string
}

// FormatFactory describes a func used for instantiating a Format during assertions.
type FormatFactory func(string) format.Format

// Name asserts the correct behaviour of github.com/corvus-ch/horcrux/format.Format.OutputFileName().
func Name(t *testing.T, tests []NameTest, factory FormatFactory) {
	for _, test := range tests {
		t.Run(fmt.Sprint(test.X), func(t *testing.T) {
			f := factory(test.Stem)
			assert.Equal(t, test.Expected, f.OutputFileName(test.X))
		})
	}
}
