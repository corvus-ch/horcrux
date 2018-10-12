package internal

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bketelsen/logr"
	"github.com/corvus-ch/horcrux/create"
	"github.com/corvus-ch/horcrux/format"
	"gopkg.in/alecthomas/kingpin.v2"
)

type createAction func(cfg create.Config) error

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
	return c.action(c)
}

func (c *createCommand) Input() (io.Reader, error) {
	if c.input == "-" || c.input == "" {
		return os.Stdin, nil
	}

	return os.Open(c.input)
}

func (c *createCommand) Formats() ([]format.Format, error) {
	formats := make([]format.Format, len(c.formats))

	stem := c.stem()
	for i, f := range c.formats {
		ff, err := format.New(f, stem)
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

func (c *createCommand) stem() string {
	if c.stemFlag != "" {
		return c.stemFlag
	}

	file := c.input

	if file == "-" || file == "" {
		file = "part"
	}

	b := filepath.Base(file)

	return strings.TrimSuffix(b, filepath.Ext(b))
}
