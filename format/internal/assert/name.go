package assert

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// NameTest holds the data required to assert the behaviour of github.com/corvus-ch/horcrux/format.Format.
type NameTest struct {
	X        byte
	Stem     string
	Expected string
}

// Name asserts the correct behaviour of github.com/corvus-ch/horcrux/format.Format.OutputFileName().
func Name(t *testing.T, tests []NameTest, factory FormatFactory) {
	for _, test := range tests {
		t.Run(fmt.Sprint(test.X), func(t *testing.T) {
			f := factory(test.Stem)
			assert.Equal(t, test.Expected, f.OutputFileName(test.X))
		})
	}
}
