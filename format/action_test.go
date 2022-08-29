package format_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/horcrux/format/raw"
	"github.com/corvus-ch/horcrux/input"
	"github.com/corvus-ch/logr/buffered"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFormat(t *testing.T) {
	dir := createDir(t)
	defer os.RemoveAll(dir)

	f := raw.New(input.NewStreamInput(outputStem(dir, raw.Name)))
	cfg := NewConfig(t.Name(), f)
	log := buffered.New(1)

	err := format.DoFormat(cfg, log)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, cfg)

	files, _ := filepath.Glob(fmt.Sprintf("%s.*", outputStem(dir, raw.Name)))

	assert.Empty(t, log.Buf())
	assert.Len(t, files, 1)
}

var errorTests = []struct {
	name   string
	setup  func(t *testing.T, cfg *Config, dir string)
	output string
}{
	{"input", func(t *testing.T, cfg *Config, _ string) {
		cfg.On("Input").Maybe().Return(nil, errors.New("input error"))
	}, "failed to open input: input error"},

	{"formats", func(t *testing.T, cfg *Config, _ string) {
		cfg.On("Input").Return(bytes.NewBufferString(t.Name()), nil)
		cfg.On("Formats").Return(nil, errors.New("format error"))
	}, "failed to setup output formatting: format error"},

	{"copy", func(t *testing.T, cfg *Config, dir string) {
		f, _ := os.Open(t.Name())
		cfg.On("Input").Return(f, nil)
		cfg.On("InputInfo").Maybe().Return(input.NewStreamInput(outputStem(dir, raw.Name)))
		cfg.On("Formats").Maybe().Return([]format.Format{raw.New(input.NewStreamInput(outputStem(dir, raw.Name)))}, nil)
	}, "failed to process data: invalid argument"},
}

func TestErrors(t *testing.T) {
	for _, test := range errorTests {
		t.Run(test.name, func(t *testing.T) {
			dir := createDir(t)
			defer os.RemoveAll(dir)
			cfg := &Config{}
			log := buffered.New(1)
			test.setup(t, cfg, dir)

			err := format.DoFormat(cfg, log)

			assert.Error(t, err, test.output)
			assert.Empty(t, log.Buf())
			mock.AssertExpectationsForObjects(t, cfg)
		})
	}
}

func createDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "horcrux_create")
	if err != nil {
		t.Error(err)
	}

	return dir
}

func outputStem(dir, format string) string {
	return fmt.Sprintf("%s/%s", dir, format)
}
