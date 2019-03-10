package main

import (
	"testing"

	"github.com/corvus-ch/horcrux/internal"
	"github.com/corvus-ch/logr/buffered"
	"github.com/sebdah/goldie"
	"github.com/stretchr/testify/assert"
)

var appTests = map[string][]string{
	"default": {"help"},
	"create":  {"help", "create"},
	"restore": {"help", "restore"},
}

func TestApp(t *testing.T) {
	for name, args := range appTests {
		t.Run(name, func(t *testing.T) {
			cfg, err := internal.NewConfigFromYaml("")
			assert.NoError(t, err)
			log := buffered.New(1)
			_, err = App(cfg, log).Terminate(nil).Parse(args)
			assert.Nil(t, err)
			goldie.Assert(t, name, log.Buf().Bytes())
		})
	}
}
