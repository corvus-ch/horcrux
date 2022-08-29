package internal

import (
	"github.com/bketelsen/logr"
	"github.com/corvus-ch/horcrux/create"
	"github.com/corvus-ch/horcrux/format"
	"github.com/corvus-ch/horcrux/restore"
	"github.com/corvus-ch/logr/writer_adapter"
	"gopkg.in/alecthomas/kingpin.v2"
)

func App(log logr.Logger) *kingpin.Application {
	w := writer_adapter.NewBufferedErrorWriter(log)
	app := kingpin.New("horcrux", "a helper for preparing backups of data worth protecting")
	app.UsageWriter(w)
	app.ErrorWriter(w)
	RegisterCreateCommand(app, log, create.Create)
	RegisterFormatCommand(app, log, format.DoFormat)
	RegisterRestoreCommand(app, log, restore.Restore)

	return app
}
