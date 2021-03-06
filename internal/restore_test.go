package internal_test

import (
	"github.com/corvus-ch/horcrux/format/text"
	"io/ioutil"
	"os"
	"testing"

	"github.com/bketelsen/logr"
	"github.com/corvus-ch/horcrux/format/base64"
	"github.com/corvus-ch/horcrux/format/raw"
	"github.com/corvus-ch/horcrux/format/zbase32"
	"github.com/corvus-ch/horcrux/internal"
	"github.com/corvus-ch/horcrux/restore"
	"github.com/corvus-ch/logr/buffered"
	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
)

func assertRestoreAction(t *testing.T, args []string, action func(restore.Config, restore.PasswordProvider, logr.Logger) error) {
	log := buffered.New(0)
	app := kingpin.New("test", "test")
	internal.RegisterRestoreCommand(app, log, action)
	_, err := app.Parse(args)
	assert.Nil(t, err)
}

func createTempFile(t *testing.T) *os.File {
	file, err := ioutil.TempFile("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	return file
}

func TestRestoreCommand_Decrypt(t *testing.T) {
	file := createTempFile(t)
	defer os.Remove(file.Name())
	tests := []struct {
		name      string
		args      []string
		encrypted bool
	}{
		{"plain", []string{"restore", file.Name()}, false},
		{"encrypted", []string{"restore", "-d", file.Name()}, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertRestoreAction(t, test.args, func(cfg restore.Config, _ restore.PasswordProvider, _ logr.Logger) error {
				assert.Equal(t, test.encrypted, cfg.Decrypt())
				return nil
			})
		})
	}
}

func TestRestoreCommand_FileNames(t *testing.T) {
	file := createTempFile(t)
	defer os.Remove(file.Name())
	tests := []struct {
		name  string
		args  []string
		names []string
	}{
		{"single", []string{"restore", file.Name()}, []string{file.Name()}},
		{"multiple", []string{"restore", file.Name(), file.Name()}, []string{file.Name(), file.Name()}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertRestoreAction(t, test.args, func(cfg restore.Config, _ restore.PasswordProvider, _ logr.Logger) error {
				assert.Equal(t, test.names, cfg.FileNames())
				return nil
			})
		})
	}
}

func TestRestoreCommand_Format(t *testing.T) {
	file := createTempFile(t)
	defer os.Remove(file.Name())
	tests := []struct {
		name   string
		args   []string
		format string
	}{
		{"default", []string{"restore", file.Name()}, text.Name},
		{"raw", []string{"restore", "-f", "raw", file.Name()}, raw.Name},
		{"zbase32", []string{"restore", "-f", "zbase32", file.Name()}, zbase32.Name},
		{"base64", []string{"restore", "-f", "base64", file.Name()}, base64.Name},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertRestoreAction(t, test.args, func(cfg restore.Config, _ restore.PasswordProvider, _ logr.Logger) error {
				format, err := cfg.Format()
				if err != nil {
					t.Fatal("Missing format")
				}
				assert.Equal(t, test.format, format.Name())
				return nil
			})
		})
	}
}

func TestRestoreCommand_Output(t *testing.T) {
	input := createTempFile(t)
	defer os.Remove(input.Name())
	output := createTempFile(t)
	defer os.Remove(input.Name())
	tests := []struct {
		name   string
		args   []string
		writer *os.File
	}{
		{"stdout", []string{"restore", input.Name()}, os.Stdout},
		{"file", []string{"restore", "-o", output.Name(), input.Name()}, output},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertRestoreAction(t, test.args, func(cfg restore.Config, _ restore.PasswordProvider, _ logr.Logger) error {
				reader, err := cfg.Output()
				assert.Nil(t, err)
				assert.Equal(t, test.writer.Name(), reader.(*os.File).Name())
				return nil
			})
		})
	}
}
