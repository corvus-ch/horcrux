package internal_test

import (
	"os"
	"testing"

	"github.com/bketelsen/logr"
	"github.com/corvus-ch/horcrux/create"
	"github.com/corvus-ch/horcrux/internal"
	"github.com/corvus-ch/logr/buffered"
	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
)

func assertCreateAction(t *testing.T, args []string, action func(create.Config, logr.Logger) error) {
	log := buffered.New(0)
	app := kingpin.New("test", "test")
	internal.RegisterCreateCommand(app, log, action)
	_, err := app.Parse(args)
	assert.Nil(t, err)
}

func TestCreateCommand_Encrypt(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		encrypted bool
	}{
		{"plain", []string{"create"}, false},
		{"encrypted", []string{"create", "-e"}, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertCreateAction(t, test.args, func(cfg create.Config, _ logr.Logger) error {
				assert.Equal(t, test.encrypted, cfg.Encrypt())
				return nil
			})
		})
	}
}

func TestCreateCommand_Threshold(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		threshold int
	}{
		{"default", []string{"create"}, 2},
		{"flagged", []string{"create", "-n", "42"}, 42},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertCreateAction(t, test.args, func(cfg create.Config, _ logr.Logger) error {
				assert.Equal(t, test.threshold, cfg.Threshold())
				return nil
			})
		})
	}
}

func TestCreateCommand_Parts(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		parts int
	}{
		{"default", []string{"create"}, 3},
		{"flagged", []string{"create", "-m", "42"}, 42},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertCreateAction(t, test.args, func(cfg create.Config, _ logr.Logger) error {
				assert.Equal(t, test.parts, cfg.Parts())
				return nil
			})
		})
	}
}

func TestCreateCommand_Input(t *testing.T) {
	file := tmpFile(t)
	defer os.Remove(file.Name())
	for name, test := range newInputTests("create", file) {
		t.Run(name, func(t *testing.T) {
			assertCreateAction(t, test.args, func(cfg create.Config, _ logr.Logger) error {
				reader, err := cfg.Input()
				assert.NoError(t, err)
				assert.Equal(t, test.file.Name(), reader.(*os.File).Name())
				return nil
			})
		})
	}
}

func TestCreateCommand_Formats(t *testing.T) {
	for name, test := range newFormatsTests("create") {
		t.Run(name, func(t *testing.T) {
			assertCreateAction(t, test.args, func(cfg create.Config, _ logr.Logger) error {
				assertFormats(t, cfg, test.formats)
				return nil
			})
		})
	}
}
