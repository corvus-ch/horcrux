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

// Parser represents the text parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// Parse parses a data line.
func (p *Parser) Parse() (*line, error) {
	ln := &line{}
	var err error

	if err = p.readNonDataLines(); err != nil {
		return nil, err
	}

	if ln.N, err = p.readLineNumber(); err != nil {
		return nil, err
	}

	if ln.Data, err = p.readData(); err != nil {
		return nil, err
	}

	if ln.CRC, err = p.readChecksum(); err != nil {
		return nil, err
	}

	return ln, nil
}

func (p *Parser) readNonDataLines() error {
	for {
		tok, _ := p.scanIgnoreWhitespace()
		if EOF == tok {
			return io.EOF
		} else if LINO == tok {
			p.unscan()
			break
		}
	}

	return nil
}

func (p *Parser) readLineNumber() (uint64, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != LINO {
		return 0, fmt.Errorf("found %q, expected line number", lit)
	}
	n, err := strconv.ParseUint(lit, 10, 64)
	if nil != err {
		return n, fmt.Errorf("faled to parse integer %q: %v", lit, err)
	}

	return n, nil
}

func (p *Parser) readData() (string, error) {
	var data bytes.Buffer
	for {
		tok, lit := p.scanIgnoreWhitespace()
		if EOF == tok || EOL == tok {
			return "", fmt.Errorf("found %q, expected data or checksum", lit)
		}
		if DESC == tok {
			return "", fmt.Errorf("found %q, expected data or checksum", lit)
		}
		if CRC == tok {
			p.unscan()
			break
		}
		data.WriteString(lit)
	}

	return data.String(), nil
}

func (p *Parser) readChecksum() (uint32, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if CRC != tok {
		return 0, fmt.Errorf("found %q, expected checksum", lit)
	}
	n, err := strconv.ParseUint(lit, 16, 32)
	if nil != err {
		return 0, fmt.Errorf("faled to parse integer %q: %v", lit, err)
	}

	return uint32(n), nil
}

// Returns the next token from the underlying Scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the Scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// Scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

// Pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }
