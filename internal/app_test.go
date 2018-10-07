package internal_test

import (
	"testing"

	"github.com/corvus-ch/horcrux/internal"
	"github.com/corvus-ch/logr/buffered"
	"github.com/sebdah/goldie"
	"github.com/stretchr/testify/assert"
)

var appTests = []struct {
	name string
	args []string
}{
	{"default", []string{"help"}},
	{"create", []string{"help", "create"}},
	{"restore", []string{"help", "restore"}},
}

func TestApp(t *testing.T) {
	for _, test := range appTests {
		t.Run(test.name, func(t *testing.T) {
			log := buffered.New(1)
			app := internal.App(log).Terminate(nil)
			_, err := app.Parse(test.args)
			assert.Nil(t, err)
			goldie.Assert(t, t.Name(), log.Buf().Bytes())
		})
	}
}
