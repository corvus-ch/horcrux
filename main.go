package main

import (
	"fmt"
	stdLog "log"
	"os"

	"github.com/corvus-ch/horcrux/internal"
	"github.com/corvus-ch/logr/std"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	log := std.New(1, stdLog.New(os.Stderr, "", 0), stdLog.New(os.Stdout, "", 0))
	app := internal.App(log)
	app.Version(fmt.Sprintf("%v, commit %v, built at %v", version, commit, date))
	_, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Errorf("%s, try --help", err)
		os.Exit(1)
	}
}
