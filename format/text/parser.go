package text

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

// Represents a data line.
type line struct {
	N    uint64 // The line number.
	Data string // The line data.
	CRC  uint32 // The CRC for the line or if data is nil for the whole document.
}

// Represents a parser.
type parser struct {
	s   *scanner
	buf struct {
		tok token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a instance of parser.
func NewParser(r io.Reader) *parser {
	return &parser{s: NewScanner(r)}
}

// Parse parses a data line.
func (p *parser) Parse() (*line, error) {
	ln := &line{}
	var tok token
	var lit string
	var n uint64
	var err error

	// Read and discard any non line data.
	for {
		tok, _ = p.scanIgnoreWhitespace()
		if EOF == tok {
			return nil, io.EOF
		} else if LINO == tok {
			p.unscan()
			break
		}
	}

	// Read the line number.
	tok, lit = p.scanIgnoreWhitespace()
	if tok != LINO {
		return nil, fmt.Errorf("found %q, expected line number", lit)
	}
	n, err = strconv.ParseUint(lit, 10, 64)
	if nil != err {
		return nil, fmt.Errorf("Faled to parse integer %q: %v", lit, err)
	}
	ln.N = n

	// Read all the line data.
	var data bytes.Buffer
	for {
		tok, lit = p.scanIgnoreWhitespace()
		if EOF == tok || EOL == tok {
			return nil, fmt.Errorf("found %q, expected data or checksum", lit)
		}
		if DESC == tok {
			return nil, fmt.Errorf("found %q, expected data or checksum", lit)
		}
		if CRC == tok {
			p.unscan()
			break
		}
		data.WriteString(lit)
	}
	ln.Data = data.String()

	// Read the lines CRC checksum.
	if tok, lit := p.scanIgnoreWhitespace(); CRC != tok {
		return nil, fmt.Errorf("found %q, expected checksum", lit)
	}
	n, err = strconv.ParseUint(lit, 16, 32)
	if nil != err {
		return nil, fmt.Errorf("Faled to parse integer %q: %v", lit, err)
	}
	ln.CRC = uint32(n)

	return ln, nil
}

// Returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *parser) scan() (tok token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// Scans the next non-whitespace token.
func (p *parser) scanIgnoreWhitespace() (tok token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

// Pushes the previously read token back onto the buffer.
func (p *parser) unscan() { p.buf.n = 1 }
