package internal

import (
	"io"
	"os"

	"github.com/bketelsen/logr"
	"github.com/corvus-ch/horcrux/create"
	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/horcrux/input"
	"gopkg.in/alecthomas/kingpin.v2"
)

type createAction func(cfg create.Config, logger logr.Logger) error

type createCommand struct {
	action createAction

	// The logger
	log logr.Logger

	// Arguments
	input string

	// Flags
	encrypt   bool
	formats   []string
	parts     int
	stemFlag  string
	threshold int

	// internal
	info input.Input
}

func RegisterCreateCommand(app *kingpin.Application, log logr.Logger, action createAction) *createCommand {
	c := &createCommand{action: action, log: log}

	cc := app.Command("create", "create a new set of horcruxes")
	cc.Action(c.Execute)
	cc.Arg("input", "the input file to split").
		StringVar(&c.input)
	cc.Flag("output", "name stem for the output files").
		Short('o').
		StringVar(&c.stemFlag)
	cc.Flag("sharecount", "the number of horcruxes to create").
		Default("3").
		Short('m').
		IntVar(&c.parts)
	cc.Flag("threshold", "the minimal number of horcruxes required for a restore").
		Default("2").
		Short('n').
		IntVar(&c.threshold)
	cc.Flag("format", "the formats the horcruxes are created in").
		Default(format.Default).
		Short('f').
		StringsVar(&c.formats)
	cc.Flag("encrypt", "encrypt output").
		Short('e').
		BoolVar(&c.encrypt)

	return c
}

func (c *createCommand) Execute(_ *kingpin.ParseContext) error {
	return c.action(c, c.log)
}

func (c *createCommand) Input() (io.Reader, error) {
	if c.input == "-" || c.input == "" {
		return os.Stdin, nil
	}

	return os.Open(c.input)
}

func (c *createCommand) InputInfo() input.Input {
	if c.info != nil {
		return c.info
	}

	if c.input == "-" || c.input == "" {
		c.info = input.NewStreamInput(c.stemFlag)
	} else {

		c.info = input.NewFileInput(c.input, c.stemFlag)
	}

	return c.info
}

func (c *createCommand) Formats() ([]format.Format, error) {
	formats := make([]format.Format, len(c.formats))

	for i, f := range c.formats {
		ff, err := format.New(f, c.InputInfo())
		if err != nil {
			return nil, err
		}
		formats[i] = ff
	}

	return formats, nil
}

func (c *createCommand) Encrypt() bool {
	return c.encrypt
}

func (c *createCommand) Parts() int {
	return c.parts
}

func (c *createCommand) Threshold() int {
	return c.threshold
}
