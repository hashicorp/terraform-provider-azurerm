package datetime

import (
	"bufio"
	"bytes"
	"io"
)

var eof = rune(0)

type token int

// Parsing tokens for internal use, only capitalized for stylistic reasons.
const (
	ILLEGAL token = iota
	EOF
	NUMBER
	DASH
	COLON
	DOT
	PLUS
	T
	Z
)

func isNumber(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

type scanner struct {
	r *bufio.Reader
}

func newScanner(r io.Reader) *scanner {
	return &scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *scanner) unread() { _ = s.r.UnreadRune() }

// scan returns the next token and literal value.
func (s *scanner) scan() (tok token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	if isNumber(ch) {
		s.unread()
		return s.scanNumber()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return EOF, ""
	case '-':
		return DASH, string(ch)
	case ':':
		return COLON, string(ch)
	case '.':
		return DOT, string(ch)
	case '+':
		return PLUS, string(ch)
	case 'T':
		return T, string(ch)
	case 'Z':
		return Z, string(ch)
	}

	return ILLEGAL, string(ch)
}

func (s *scanner) scanNumber() (tok token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isNumber(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return NUMBER, buf.String()
}
