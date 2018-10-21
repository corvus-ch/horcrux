package base64_test

import (
	"testing"

	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/horcrux/format/base64"
	formatAssert "github.com/corvus-ch/horcrux/format/internal/assert"
	"github.com/stretchr/testify/assert"
)

func factory(s string) format.Format {
	return base64.New(s)
}

func TestFormat_Reader(t *testing.T) {
	formatAssert.DataRead(t, factory, ".base64")
}

func TestFormat_Writer(t *testing.T) {
	formatAssert.DataWrite(t, factory, ".base64")
}

func TestFormat_Name(t *testing.T) {
	assert.Equal(t, base64.Name, base64.New("").Name())
}
