package assert

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

// DataTest holds data required for format data read and write assertions.
type DataTest struct {
	Decoded []byte
	Encoded string
}

// DataRead asserts a formats read behaviour.
func DataRead(t *testing.T, tests []DataTest, factory FormatFactory) {
	for _, test := range tests {
		t.Run(test.Encoded, func(t *testing.T) {
			r, err := factory("").Reader(bytes.NewBufferString(test.Encoded))
			assert.Nil(t, err)
			out, err := ioutil.ReadAll(r)
			assert.NoError(t, err)
			assert.Equal(t, test.Decoded, out)
		})
	}
}

// DataWrite asserts a formats write behaviour.
func DataWrite(t *testing.T, tests []DataTest, factory FormatFactory) {
	for _, test := range tests {
		t.Run(test.Encoded, func(t *testing.T) {
			out := &bytes.Buffer{}
			w, cl, err := factory("").Writer(out)
			assert.Len(t, cl, 1)
			assert.NoError(t, err)
			w.Write(test.Decoded)
			cl[0].Close()
			assert.Equal(t, test.Encoded, out.String())
		})
	}
}
