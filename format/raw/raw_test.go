package raw_test

import (
	"testing"

	"github.com/corvus-ch/horcrux/format"
	formatAssert "github.com/corvus-ch/horcrux/format/internal/assert"
	"github.com/corvus-ch/horcrux/format/raw"
	"github.com/stretchr/testify/assert"
)

func factory(s string) format.Format {
	return raw.New(s)
}

func TestFormat_Reader(t *testing.T) {
	formatAssert.DataRead(t, factory, ".raw")
}

func TestFormat_Writer(t *testing.T) {
	formatAssert.DataWrite(t, factory, ".raw")
}

func TestFormat_Name(t *testing.T) {
	assert.Equal(t, raw.Name, raw.New("").Name())
}
