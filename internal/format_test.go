package internal_test

import (
	"os"
	"testing"

	"github.com/bketelsen/logr"
	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/horcrux/internal"
	"github.com/corvus-ch/logr/buffered"
	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
)

func assertFormatAction(t *testing.T, args []string, action func(format.Config, logr.Logger) error) {
	log := buffered.New(0)
	app := kingpin.New("test", "test")
	internal.RegisterFormatCommand(app, log, action)
	_, err := app.Parse(args)
	assert.Nil(t, err)
}

func TestFormatCommand_Input(t *testing.T) {
	file := tmpFile(t)
	defer os.Remove(file.Name())
	for name, test := range newInputTests("format", file) {
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

func TestFormatCommand_Formats(t *testing.T) {
	for name, test := range newFormatsTests("format") {
		t.Run(name, func(t *testing.T) {
			assertFormatAction(t, test.args, func(cfg format.Config, _ logr.Logger) error {
				assertFormats(t, cfg, test.formats)
				return nil
			})
		})
	}
}
