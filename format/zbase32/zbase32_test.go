package zbase32_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/corvus-ch/horcrux/format"
	formatAssert "github.com/corvus-ch/horcrux/format/internal/assert"
	"github.com/corvus-ch/horcrux/format/zbase32"
	"github.com/corvus-ch/horcrux/input"
	"github.com/stretchr/testify/assert"
)

func factory(i input.Input) format.Format {
	return zbase32.New(i)
}

func TestFormat_Reader(t *testing.T) {
	formatAssert.DataRead(t, factory, ".zbase32")
}

func TestFormat_Writer(t *testing.T) {
	formatAssert.DataWrite(t, factory, ".zbase32", func(file string, x byte) []string {
		base := filepath.Base(file)
		name := base[0 : len(base)-len(filepath.Ext(base))]
		return []string{fmt.Sprintf("%s.zbase32.%03d", name, x)}
	})
}

func TestFormat_Name(t *testing.T) {
	assert.Equal(t, zbase32.Name, zbase32.New(nil).Name())
}
