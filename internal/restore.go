package internal

import (
	"io"
	"os"

	"github.com/bketelsen/logr"
	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/horcrux/restore"
	"gopkg.in/alecthomas/kingpin.v2"
)

type restoreAction func(cfg restore.Config) error

type restoreCommand struct {
	action restoreAction

	// The logger
	log logr.Logger

	// Args
	files []string

	// Flags
	decrypt bool
	format  string
	output  string
}

func RegisterRestoreCommand(app *kingpin.Application, log logr.Logger, action restoreAction) *restoreCommand {
	c := &restoreCommand{action: action, log: log}

	cc := app.Command("restore", "restores your valuable data from a set of horcruxes")
	cc.Action(c.Execute)
	cc.Arg("files", "path tho the individual horcruxes").
		Required().
		ExistingFilesVar(&c.files)
	cc.Flag("output", "path to the output").
		Short('o').
		StringVar(&c.output)
	cc.Flag("format", "the formats the horcruxes are created in").
		Default(format.DEFAULT).
		Short('f').
		StringVar(&c.format)
	cc.Flag("decrypt", "encrypt output").
		Short('d').
		BoolVar(&c.decrypt)

	return c
}

func (c *restoreCommand) Execute(_ *kingpin.ParseContext) error {
	return c.action(c)
}

func (c *restoreCommand) Format() (format.Format, error) {
	return format.New(c.format, "")
}

func (c *restoreCommand) Decrypt() bool {
	return c.decrypt
}

func (c *restoreCommand) Output() (io.Writer, error) {
	if len(c.output) == 0 || c.output == "-" {
		return os.Stdout, nil
	}

	return os.Create(c.output)
}

func (c *restoreCommand) FileNames() []string {
	return c.files
}
