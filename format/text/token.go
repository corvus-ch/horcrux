package text

// token represents a lexical token.
type token int

const (
	// Special tokens
	EOF token = iota // End of file
	EOL              // End of line
	WS               // White space

	// Literals
	CRC  // Checksum
	DATA // Payload data
	DESC // Arbitrary description
	LINO // Line number
)
