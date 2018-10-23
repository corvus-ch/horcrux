package create_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	action "github.com/corvus-ch/horcrux/create"
	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/horcrux/format/raw"
	"github.com/corvus-ch/horcrux/meta"
	"github.com/corvus-ch/logr/buffered"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var filePattern = regexp.MustCompile(`(.*)\.[^[\.]+\.(\d{3})$`)

var createTests = []struct {
	name      string
	encrypted bool
}{
	{"raw", false},
	{"encrypted", true},
}

func TestCreate(t *testing.T) {
	for _, test := range createTests {
		t.Run(test.name, func(t *testing.T) {
			dir := createDir(t)
			defer os.RemoveAll(dir)

			f := raw.New(meta.NewInputMock(outputStem(dir, raw.Name)))
			cfg := NewConfig(t.Name(), f, test.encrypted)
			log := buffered.New(1)

			err := action.Create(cfg, log)

			assert.Nil(t, err)
			mock.AssertExpectationsForObjects(t, cfg)

			files, _ := filepath.Glob(fmt.Sprintf("%s.*", outputStem(dir, raw.Name)))

			if test.encrypted {
				files = decrypt(t, files, readPasswords(t, log.Buf()))
			} else {
				assert.Empty(t, log.Buf())
			}

			assert.Len(t, files, cfg.Parts())
			assertCombine(t, t.Name(), cfg, files)
		})
	}
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

	{"processing pipeline", func(t *testing.T, cfg *Config, dir string) {
		cfg.On("Input").Return(bytes.NewBufferString(t.Name()), nil)
		cfg.On("Formats").Maybe().Return([]format.Format{raw.New(meta.NewInputMock(outputStem(dir, raw.Name)))}, nil)
		cfg.On("Parts").Return(3)
		cfg.On("Threshold").Return(4)
		cfg.On("Encrypt").Return(false)
	}, "failed to create processing pipeline: parts cannot be less than threshold"},

	{"copy", func(t *testing.T, cfg *Config, dir string) {
		f, _ := os.Open(t.Name())
		cfg.On("Input").Return(f, nil)
		cfg.On("Formats").Maybe().Return([]format.Format{raw.New(meta.NewInputMock(outputStem(dir, raw.Name)))}, nil)
		cfg.On("Parts").Return(3)
		cfg.On("Threshold").Return(2)
		cfg.On("Encrypt").Return(false)
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

			err := action.Create(cfg, log)

			assert.Error(t, err, test.output)
			assert.Empty(t, log.Buf())
			mock.AssertExpectationsForObjects(t, cfg)
		})
	}
}

func assertCombine(t *testing.T, expect string, cfg action.Config, files []string) {
	buf := &bytes.Buffer{}
	args := append([]string{"-o", "-"}, files[cfg.Parts()-cfg.Threshold():]...)
	combine := exec.Command("gfcombine", args...)
	combine.Stdout = buf
	combine.Stderr = buf
	assert.Nil(t, combine.Run())
	assert.Equal(t, expect, buf.String())
}

func readPasswords(t *testing.T, buf *bytes.Buffer) map[string]string {
	passwords := make(map[string]string, 3)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			if line == "" {
				break
			}
		}

		parts := strings.SplitN(line[18:], ": ", 2)
		passwords[parts[0]] = parts[1]
	}

	return passwords
}

func decrypt(t *testing.T, files []string, passwords map[string]string) []string {
	for i, f := range files {
		of := filePattern.ReplaceAllString(f, `$1.gpg.$2`)
		decrypt := exec.Command(
			"gpg",
			"--batch",
			"--passphrase-fd", "0",
			"--decrypt",
			"--output", of,
			f,
		)
		index := f[len(f)-3:]
		password, ok := passwords[index]
		if !ok {
			t.Errorf("No password found for index %s", index)
		}
		decrypt.Stdin = bytes.NewBufferString(password)
		assert.Nil(t, decrypt.Run())
		files[i] = of
	}

	return files
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
