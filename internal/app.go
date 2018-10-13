package internal

import (
	"github.com/bketelsen/logr"
	"github.com/corvus-ch/logr/writer_adapter"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/corvus-ch/horcrux/create"
)

func App(log logr.Logger) *kingpin.Application {
	w := writer_adapter.NewBufferedErrorWriter(log)
	app := kingpin.New("horcrux", "paper backup for the paranoid")
	app.UsageWriter(w)
	app.ErrorWriter(w)
	RegisterCreateCommand(app, log, create.Create)
	RegisterRestoreCommand(app, log, nil)

	return app
}
