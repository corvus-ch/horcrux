package text

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Ensure the Parser can parse strings into line ASTs.
func TestParserParse(t *testing.T) {
	var tests = []struct {
		name string
		s    string
		line *line
		err  string
	}{
		{
			name: "full length line",
			s:    `1: bta4s yx864 9yjku rxzwk zm886 giua8 adgo3 8go4u F34870`,
			line: &line{1, "bta4syx8649yjkurxzwkzm886giua8adgo38go4u", 0xF34870},
		},

		{
			name: "partial line",
			s:    `6: j7n7s by 9C8FC7`,
			line: &line{6, "j7n7sby", 0x9C8FC7},
		},

		{
			name: "final checksum line",
			s:    `7: A78E92`,
			line: &line{7, ``, 0xA78E92},
		},

		{
			name: "with leading description",
			s:    "Integer posuere erat a ante venenatis dapibus posuere velit aliquet.\n1: j7n7s by 9C8FC7",
			line: &line{1, "j7n7sby", 0x9C8FC7},
		},

		// Errors
		{name: "eof", s: `foo`, err: io.EOF.Error()},
		{name: "invalid data", s: `1: Lorem impsum`, err: `found "Lorem", expected data or checksum`},
		{name: "missing checksum", s: `6: j7n7s by`, err: `found "", expected data or checksum`},
		{name: "none zbase32 chars", s: `42: j7n7s lorem`, err: `found "lorem", expected data or checksum`},
		{name: "multi line", s: "6 j7n7s by\n7: as8kg", err: `found "", expected data or checksum`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ln, err := NewParser(strings.NewReader(tt.s)).Parse()
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
			}
			assert.Equal(t, tt.line, ln)
		})
	}
}
