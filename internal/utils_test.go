package internal_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bketelsen/logr"
	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/horcrux/format/raw"
	"github.com/corvus-ch/horcrux/format/text"
	"github.com/corvus-ch/horcrux/format/zbase32"
	"github.com/stretchr/testify/assert"
)

type inputTests map[string]struct {
	args []string
	file *os.File
}

type formatTests map[string]struct {
	args    []string
	formats []string
}

func newInputTests(action string, file *os.File) inputTests {
	return inputTests{
		"default": {[]string{action}, os.Stdin},
		"stdin":   {[]string{action, "--", "-"}, os.Stdin},
		"file":    {[]string{action, "--", file.Name()}, file},
	}
}

func newFormatsTests(action string) formatTests {
	return formatTests{
		"default":  {[]string{action}, []string{text.Name}},
		"single":   {[]string{action, "-f", "raw"}, []string{raw.Name}},
		"multiple": {[]string{action, "-f", "raw", "-f", "zbase32"}, []string{raw.Name, zbase32.Name}},
	}
}

func tmpFile(t *testing.T) *os.File {
	file, err := ioutil.TempFile("", t.Name())
	if err != nil {
		t.Fatal(err)
	}

	return file
}

func assertFormats(t *testing.T, cfg format.Config, testFormats []string) {
	formats, err := cfg.Formats()
	assert.NoError(t, err)
	assert.Equal(t, len(testFormats), len(formats))
	for i, name := range testFormats {
		assert.Equal(t, name, formats[i].Name())
	}
}

func inputTest(t *testing.T, action string) {
	file := tmpFile(t)
	defer os.Remove(file.Name())
	for name, test := range newInputTests(action, file) {
		t.Run(name, func(t *testing.T) {
			assertFormatAction(t, test.args, func(cfg format.Config, _ logr.Logger) error {
				reader, err := cfg.Input()
				assert.NoError(t, err)
				assert.Equal(t, test.file.Name(), reader.(*os.File).Name())
				return nil
			})
		})
	}
}

func formatTest(t *testing.T, action string) {
	for name, test := range newFormatsTests(action) {
		t.Run(name, func(t *testing.T) {
			assertFormatAction(t, test.args, func(cfg format.Config, _ logr.Logger) error {
				assertFormats(t, cfg, test.formats)
				return nil
			})
		})
	}
}
