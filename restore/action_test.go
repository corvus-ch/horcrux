package restore_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/corvus-ch/horcrux/format/raw"
	action "github.com/corvus-ch/horcrux/restore"
	"github.com/corvus-ch/logr/buffered"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var restoreTests = []struct {
	name      string
	encrypted bool
}{
	{"raw", false},
	{"encrypted", true},
}

func TestRestore(t *testing.T) {
	for _, test := range restoreTests {
		t.Run(test.name, func(t *testing.T) {
			dir := createDir(t)
			defer os.RemoveAll(dir)
			files := split(t, dir, raw.Name, t.Name(), test.encrypted)

			cfg, buf := NewConfig(files, raw.New("raw"), test.encrypted)
			prompt := NewPrompter(files)
			log := buffered.New(1)

			err := action.Restore(cfg, prompt, log)

			assert.Nil(t, err)
			mock.AssertExpectationsForObjects(t, cfg)
			assert.Equal(t, t.Name(), buf.String())
		})
	}
}

var errorTests = []struct {
	name      string
	encrypted bool
	setup     func(t *testing.T, cfg *Config, dir string, files []string) (*bytes.Buffer, string, []string)
	output    string
}{
	{"index parsing", false, func(t *testing.T, cfg *Config, dir string, files []string) (*bytes.Buffer, string, []string) {
		buf := &bytes.Buffer{}
		cfg.On("Output").Maybe().Return(buf, nil)
		cfg.On("FileNames").Return(files)
		f, err := os.Create(fmt.Sprintf("%s/data.%s.abc", dir, raw.Name))
		if err != nil {
			t.Error(err)
		}
		defer f.Close()
		files[0] = f.Name()

		return &bytes.Buffer{}, fmt.Sprintf("INFO Reading file %s/data.raw.abc\n", dir), files
	}, "failed to get shamir index from file name: strconv.ParseUint: parsing \"abc\": invalid syntax"},

	{"file open", false, func(t *testing.T, cfg *Config, dir string, files []string) (*bytes.Buffer, string, []string) {
		buf := &bytes.Buffer{}
		cfg.On("Output").Maybe().Return(buf, nil)
		cfg.On("FileNames").Return(files)
		files[0] = fmt.Sprintf("%s/data.%s.000", dir, raw.Name)

		return buf, fmt.Sprintf("INFO Reading file %[1]s/data.raw.000\n", dir), files
	}, "failed to open {{ dir }}/data.raw.000: open {{ dir }}/data.raw.000: no such file or directory"},

	{"formats", false, func(t *testing.T, cfg *Config, dir string, files []string) (*bytes.Buffer, string, []string) {
		buf := &bytes.Buffer{}
		cfg.On("Output").Maybe().Return(buf, nil)
		cfg.On("FileNames").Return(files)
		cfg.On("Format").Return(nil, errors.New("format error"))

		return buf, fmt.Sprintf("INFO Reading file %s\n", files[0]), files
	}, "failed to to detect input format: format error"},

	{"missing parts", false, func(t *testing.T, cfg *Config, dir string, files []string) (*bytes.Buffer, string, []string) {
		buf := &bytes.Buffer{}
		cfg.On("Output").Maybe().Return(buf, nil)
		cfg.On("FileNames").Return(files[2:])
		cfg.On("Format").Return(raw.New("stem"), nil)
		cfg.On("Decrypt").Return(false)

		return buf, fmt.Sprintf("INFO Reading file %s\n", files[2]), files
	}, "failed to create processing pipeline: at least two parts are required to reconstruct the secret"},

	{"output", false, func(t *testing.T, cfg *Config, dir string, files []string) (*bytes.Buffer, string, []string) {
		buf := &bytes.Buffer{}
		cfg.On("Output").Maybe().Return(nil, errors.New("output error"))
		cfg.On("FileNames").Return(files[1:])
		cfg.On("Format").Return(raw.New("stem"), nil)
		cfg.On("Decrypt").Return(false)

		return buf, fmt.Sprintf("INFO Reading file %s\nINFO Reading file %s\n", files[1], files[2]), files
	}, "failed to open output: output error"},

	{"decryption", true, func(t *testing.T, cfg *Config, dir string, files []string) (*bytes.Buffer, string, []string) {
		buf := &bytes.Buffer{}
		cfg.On("Output").Maybe().Return(buf, nil)
		cfg.On("FileNames").Return(files)
		cfg.On("Format").Return(raw.New("stem"), nil)
		cfg.On("Decrypt").Return(true)

		return buf, fmt.Sprintf("INFO Reading file %s\n", files[0]), nil
	}, "failed to create decryption reader: no passwords available"},
}

func TestRestoreErrors(t *testing.T) {
	for _, test := range errorTests {
		t.Run(test.name, func(t *testing.T) {
			dir := createDir(t)
			defer os.RemoveAll(dir)
			files := split(t, dir, raw.Name, t.Name(), test.encrypted)
			cfg := &Config{}
			buf, output, passwords := test.setup(t, cfg, dir, files)
			prompt := NewPrompter(passwords)
			log := buffered.New(0)

			err := action.Restore(cfg, prompt, log)

			assert.Error(t, err)
			assert.Equal(t, strings.Replace(test.output, "{{ dir }}", dir, -1), err.Error())
			mock.AssertExpectationsForObjects(t, cfg)
			assert.Empty(t, buf.String())
			assert.Equal(t, output, log.Buf().String())
		})
	}
}

func split(t *testing.T, dir, format, data string, encrypted bool) []string {
	file, err := os.Create(fmt.Sprintf("%s/data.%s", dir, format))
	if err != nil {
		t.Error(err)
	}
	if _, err := file.WriteString(data); err != nil {
		t.Error(err)
	}

	split := exec.Command("gfsplit", "-n", "2", "-m", "3", file.Name())
	split.Dir = dir
	if err := split.Run(); err != nil {
		t.Error(err)
	}

	files, err := filepath.Glob(fmt.Sprintf("%s/data.%s.*", dir, format))
	if err != nil {
		t.Error(err)
	}

	if encrypted {
		for _, file := range files {
			output := strings.Replace(file, "data."+format, "data.gpg", -1)
			encrypt := exec.Command("gpg",
				"--batch",
				"--passphrase-fd", "0",
				"--symmetric",
				"--output", output,
				file,
			)
			encrypt.Stdin = bytes.NewBufferString(output)
			if err := encrypt.Run(); err != nil {
				t.Error(err)
			}
		}

		files, err = filepath.Glob(fmt.Sprintf("%s/data.gpg.*", dir))
		if err != nil {
			t.Error(err)
		}
	}

	return files
}

func createDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "horcrux_restore")
	if err != nil {
		t.Error(err)
	}

	return dir
}
