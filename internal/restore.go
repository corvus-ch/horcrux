package internal

import (
	"github.com/bketelsen/logr"
	"github.com/corvus-ch/horcrux/restore"
	"github.com/corvus-ch/logr/writer_adapter"
	"gopkg.in/alecthomas/kingpin.v2"
)

type restoreAction func(restore.Config, restore.PasswordProvider, logr.Logger) error

func RegisterRestoreCommand(app *kingpin.Application, cfg *Config, log logr.Logger, action restoreAction) {
	c := app.Command("restore", "restores your valuable data from a set of horcruxes")
	c.Action(func(_ *kingpin.ParseContext) error {
		return action(cfg, restore.NewPasswordProvider(writer_adapter.NewInfoWriter(log)), log)
	})

	c.Arg("files", "path tho the individual horcruxes").
		Required().
		SetValue(cfg.StringsValue("files"))

	c.Flag("output", "path to the output").
		Short('o').
		SetValue(cfg.StringValue("output"))

	f := cfg.EnumValue("format", "base64", "qr", "raw", "text", "zbase32")
	c.Flag("format", "the formats the horcruxes are created in").
		Default(f.String()).
		Short('f').
		SetValue(f)

	c.Flag("decrypt", "encrypt output").
		Short('d').
		SetValue(cfg.BoolValue("encrypted"))
}
