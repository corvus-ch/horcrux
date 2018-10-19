package text

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

const (
	data = "ybndrfg8ejkmcpqxot1uwisza345h769" // Valid data characters
	crc  = "0123456789ABCDEF"                 // Valid crc characters
)

// Scanner represents the lexical text scanner.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (tok token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	// If we see a digit then consume as a number.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return EOF, ""
	case 0x0A:
		return EOL, "\n"
	}

	s.unread()
	return s.scanIdent()

}

// Consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (tok token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// Consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (tok token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof || ch == eol {
			s.unread()
			break
		} else if isWhitespace(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	str := buf.String()

	if ch := s.read(); (ch == eof || ch == eol) && isCRC(str) {
		return CRC, str
	}

	if strings.ContainsRune(str, ':') && isDigit(strings.TrimRight(str, `:`)) {
		str = strings.TrimRight(str, `:`)
		str = strings.TrimLeft(str, `0`)
		return LINO, str
	}
	if isData(str) {
		return DATA, str
	}

	return DESC, str
}

// Reads the next rune from the buffered reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// Places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// Checks if the rune is a whitespace.
func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' }

// Checks if string contains only digit characters.
func isDigit(s string) bool {
	for _, r := range s {
		if !(r >= '0' && r <= '9') {
			return false
		}

	}
	return true
}

// Checks if string contains only valid data characters.
func isData(s string) bool {
	for _, r := range s {
		if !strings.ContainsRune(data, r) {
			return false
		}
	}

	return true
}

// Checks if string contains only valid CRC characters.
func isCRC(s string) bool {
	for _, r := range s {
		if !strings.ContainsRune(crc, r) {
			return false
		}
	}

	return true
}

// Represents a marker rune for the end of the reader.
var eof = rune(0)

// Represents a marker for the end of a Line.
var eol = rune(0x0a)
