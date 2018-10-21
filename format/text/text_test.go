package text_test

import (
	"testing"

	"github.com/corvus-ch/horcrux/format"
	formatAssert "github.com/corvus-ch/horcrux/format/internal/assert"
	"github.com/corvus-ch/horcrux/format/text"
	"github.com/stretchr/testify/assert"
)

func factory(s string) format.Format {
	return text.New(s)
}

func TestFormat_Reader(t *testing.T) {
	formatAssert.DataRead(t, factory, ".txt")
}

func TestFormat_Writer(t *testing.T) {
	formatAssert.DataWrite(t, factory, ".txt")
}

func TestFormat_Name(t *testing.T) {
	assert.Equal(t, text.Name, text.New("").Name())
}
