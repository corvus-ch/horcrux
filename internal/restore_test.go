package internal_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/corvus-ch/horcrux/internal"
	"github.com/corvus-ch/horcrux/restore"
	"github.com/corvus-ch/logr/buffered"
	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
)

func assertRestoreAction(t *testing.T, args []string, action func(restore.Config) error) {
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
			assertRestoreAction(t, test.args, func(cfg restore.Config) error {
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
			assertRestoreAction(t, test.args, func(cfg restore.Config) error {
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
		name string
		args []string
	}{
		{"default", []string{"restore", file.Name()}},
		{"raw", []string{"restore", "-f", "zbase32", file.Name()}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertRestoreAction(t, test.args, func(cfg restore.Config) error {
				if _, err := cfg.Format(); err != nil {
					t.Skip("Formats not yet implemented")
				}
				t.Fatal("Formats inplemented but test not yet completed")
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
			assertRestoreAction(t, test.args, func(cfg restore.Config) error {
				reader, err := cfg.Output()
				assert.Nil(t, err)
				assert.Equal(t, test.writer.Name(), reader.(*os.File).Name())
				return nil
			})
		})
	}
}
