package main

import (
	stdLog "log"
	"os"

	"github.com/corvus-ch/horcrux/internal"
	"github.com/corvus-ch/logr/std"
)

func main() {
	log := std.New(1, stdLog.New(os.Stderr, "", 0), stdLog.New(os.Stdout, "", 0))
	_, err := internal.App(log).Parse(os.Args[1:])
	if err != nil {
		log.Errorf("%s, try --help", err)
		os.Exit(1)
	}
}
