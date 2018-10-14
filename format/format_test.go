package format_test

import (
	"fmt"
	"testing"

	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/horcrux/format/raw"
	"github.com/corvus-ch/horcrux/format/zbase32"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{raw.Name, nil},
		{zbase32.Name, nil},
		{"foo", fmt.Errorf("unknown format foo")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, err := format.New(test.name, "")
			assert.Equal(t, test.err, err)
			if f != nil {
				assert.Equal(t, test.name, f.Name())
			}
		})
	}
}
