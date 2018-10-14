package text

import (
	"strings"
	"testing"
)

// Ensures the scanner can scan tokens correctly.
func TestScannerScan(t *testing.T) {
	var tests = []struct {
		s   string
		tok token
		lit string
	}{
		// Special tokens
		{``, EOF, ""},
		{` `, WS, " "},
		{"\t", WS, "\t"},
		{"\n", EOL, "\n"},

		// Literals
		{`#`, DESC, `#`},
		{`Lorem`, DESC, `Lorem`},
		{`Zx12_3U_-`, DESC, `Zx12_3U_-`},
		{`bta4s`, DATA, `bta4s`},
		{`a65s4 f5xsk`, DATA, `a65s4`},
		{`a65s4f5xsk`, DATA, `a65s4f5xsk`},
		{`5xsk 429DB8`, DATA, `5xsk`},
		{`429DB8`, CRC, `429DB8`},
		{`3:`, LINO, `3`},
		{`07:`, LINO, `7`},
		{`013:`, LINO, `13`},
		{`42:`, LINO, `42`},
		{`1234:`, LINO, `1234`},
		{`12:34`, DESC, `12:34`},
		{`00:34`, DESC, `00:34`},

		// Regressions
		{`111948`, CRC, `111948`},
		{`6zpz9 111948`, DATA, `6zpz9`},
	}

	for i, tt := range tests {
		s := NewScanner(strings.NewReader(tt.s))
		tok, lit := s.Scan()
		if tt.tok != tok {
			t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tt.s, tt.tok, tok, lit)
		} else if tt.lit != lit {
			t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, tt.s, tt.lit, lit)
		}
	}
}
