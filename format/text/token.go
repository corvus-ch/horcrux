package text

// token represents a lexical token.
type token int

const (
	// Special tokens

	// EOF denotes the end of a file.
	EOF token = iota
	// EOL denotes the end of a Line.
	EOL
	// WS denotes any whitespace except EOF and EOL.
	WS

	// Literals

	// CRC denotes a crc24 checksum.
	CRC
	// DATA denotes the zbase32 encoded payload data.
	DATA
	// LINO denotes a Line number.
	LINO
	// DESC denotes any other type of text.
	DESC
)
