package main

import (
	"fmt"
	stdLog "log"
	"os"

	"github.com/bketelsen/logr"
	"github.com/corvus-ch/horcrux/create"
	"github.com/corvus-ch/horcrux/internal"
	"github.com/corvus-ch/horcrux/restore"
	"github.com/corvus-ch/logr/std"
	"github.com/corvus-ch/logr/writer_adapter"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func App(cfg *internal.Config, log logr.Logger) *kingpin.Application {
	w := writer_adapter.NewBufferedErrorWriter(log)
	app := kingpin.New("horcrux", "a helper for preparing backups of data worth protecting")
	app.UsageWriter(w)
	app.ErrorWriter(w)
	app.Version(fmt.Sprintf("%v, commit %v, built at %v", version, commit, date))

	internal.RegisterCreateCommand(app, cfg, log, create.Create)
	internal.RegisterRestoreCommand(app, cfg, log, restore.Restore)

	return app
}

func main() {
	cfg, err := internal.NewConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read configuration: %v\n", err)
		os.Exit(1)
	}

	log := std.New(1, stdLog.New(os.Stderr, "", 0), stdLog.New(os.Stdout, "", 0))
	_, err = App(cfg, log).Parse(os.Args[1:])
	if err != nil {
		log.Errorf("%s, try --help", err)
		os.Exit(1)
	}
}
