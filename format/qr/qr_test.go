package qr_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	code "github.com/boombuler/barcode/qr"
	"github.com/corvus-ch/horcrux/format"
	formatAssert "github.com/corvus-ch/horcrux/format/internal/assert"
	"github.com/corvus-ch/horcrux/format/qr"
	"github.com/corvus-ch/horcrux/input"
	"github.com/stretchr/testify/assert"
)

func factory(i input.Input) format.Format {
	return qr.New(i)
}

func TestFormat_Reader(t *testing.T) {
	_, err := factory(nil).Reader(bytes.NewReader([]byte{}))
	assert.Error(t, err)
}

func TestFormat_Writer(t *testing.T) {
	formatAssert.DataWrite(t, factory, ".png", func(file string, x byte) []string {
		base := filepath.Base(file)
		name := base[0 : len(base)-len(filepath.Ext(base))]
		f, err := os.Stat(file)
		if err != nil {
			return []string{}
		}
		num := qr.NumChunks(qr.Capacity(code.M), f.Size())
		names := make([]string, num)
		for i := 0; i < num; i++ {
			names[i] = fmt.Sprintf("%s.%03d.%d.png", name, x, i+1)
		}

		return names
	})
}

func TestFormat_Name(t *testing.T) {
	assert.Equal(t, qr.Name, qr.New(nil).Name())
}
