package internal

import (
	"github.com/bketelsen/logr"
	"github.com/corvus-ch/horcrux/create"
	"gopkg.in/alecthomas/kingpin.v2"
)

type createAction func(cfg create.Config, logger logr.Logger) error

func RegisterCreateCommand(app *kingpin.Application, cfg *Config, log logr.Logger, action createAction) {
	c := app.Command("create", "create a new set of horcruxes")

	c.Action(func(_ *kingpin.ParseContext) error {
		return action(cfg, log)
	})

	c.Arg("input", "the input file to split").
		SetValue(cfg.StringValue("input"))

	c.Flag("output", "name stem for the output files").
		Short('o').
		SetValue(cfg.StringValue("stem"))

	m := cfg.IntValue("parts")
	c.Flag("sharecount", "the number of horcruxes to create").
		Default(m.String()).
		Short('m').
		SetValue(m)

	n := cfg.IntValue("threshold")
	c.Flag("threshold", "the minimal number of horcruxes required for a restore").
		Default(n.String()).
		Short('n').
		SetValue(n)

	c.Flag("format", "the formats the horcruxes are created in").
		Default("text").
		Short('f').
		SetValue(cfg.EnumsValue("formats", "base64", "qr", "raw", "text", "zbase32"))

	c.Flag("encrypt", "encrypt output").
		Short('e').
		SetValue(cfg.BoolValue("encrypted"))
}
