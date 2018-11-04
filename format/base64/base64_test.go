package base64_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/horcrux/format/base64"
	formatAssert "github.com/corvus-ch/horcrux/format/internal/assert"
	"github.com/corvus-ch/horcrux/meta"
	"github.com/stretchr/testify/assert"
)

func factory(i meta.Input) format.Format {
	return base64.New(i)
}

func TestFormat_Reader(t *testing.T) {
	formatAssert.DataRead(t, factory, ".base64")
}

func TestFormat_Writer(t *testing.T) {
	formatAssert.DataWrite(t, factory, ".base64", func(file string, x byte) []string {
		base := filepath.Base(file)
		name := base[0 : len(base)-len(filepath.Ext(base))]
		return []string{fmt.Sprintf("%s.base64.%03d", name, x)}
	})
}

func TestFormat_Name(t *testing.T) {
	assert.Equal(t, base64.Name, base64.New(nil).Name())
}
