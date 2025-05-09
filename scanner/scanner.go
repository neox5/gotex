package scanner

import (
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/neox5/gotex/token"
)

// ErrorHandler is used to handle scanner errors
type ErrorHandler func(pos token.Position, msg string)

// Error describes a scanner error
type Error struct {
	Pos token.Position
	Msg string
}

// Error implements the error interface
func (e Error) Error() string {
	if e.Pos.Filename != "" || e.Pos.Line > 0 {
		return fmt.Sprintf("%s: %s", e.Pos, e.Msg)
	}
	return e.Msg
}

// PrintError returns an ErrorHandler that prints errors to w
func PrintError(w io.Writer) ErrorHandler {
	return func(pos token.Position, msg string) {
		fmt.Fprintf(w, "%s: %s\n", pos, msg)
	}
}

// Scanner structure to hold scanner state
type Scanner struct {
	// Source
	src []byte // source content

	// Positioning
	file     *token.File    // source file handle
	fset     *token.FileSet // file set for position information
	offset   int            // current offset in src
	rdOffset int            // reading offset (position after current character)
	ch       rune           // current character

	// Error handling
	errHandler ErrorHandler
}

// Init initializes or re-initializes a Scanner with a new source
func (s *Scanner) Init(fset *token.FileSet, file *token.File, src []byte, errHandler ErrorHandler) {
	s.fset = fset
	s.file = file
	s.src = src
	s.errHandler = errHandler

	s.offset = 0
	s.rdOffset = 0

	// Initialize by reading the first character
	s.next()
}

const (
	bom = 0xFEFF // byte order mark, only permitted as very first character
	eof = -1     // end of file
)

// next reads the next Unicode character into s.ch and updates positioning
func (s *Scanner) next() {
	if s.rdOffset < len(s.src) {
		s.offset = s.rdOffset
		if s.ch == '\n' {
			s.file.AddLine(s.offset)
		}
		r, w := rune(s.src[s.rdOffset]), 1
		switch {
		case r == 0:
			s.error(s.offset, "illegal character NUL")
		case r >= utf8.RuneSelf:
			// not ASCII
			r, w = utf8.DecodeRune(s.src[s.rdOffset:])
			if r == utf8.RuneError && w == 1 {
				in := s.src[s.rdOffset:]
				if s.offset == 0 &&
					len(in) >= 2 &&
					(in[0] == 0xFF && in[1] == 0xFE || in[0] == 0xFE && in[1] == 0xFF) {
					// U+FEFF BOM at start of file, encoded as big- or little-endian
					// UCS-2 (i.e. 2-byte UTF-16). Give specific error (go.dev/issue/71950).
					s.error(s.offset, "illegal UTF-8 encoding (got UTF-16)")
					s.rdOffset += len(in) // consume all input to avoid error cascade
				} else {
					s.error(s.offset, "illegal UTF-8 encoding")
				}
			} else if r == bom && s.offset > 0 {
				s.error(s.offset, "illegal byte order mark")
			}
		}
		s.rdOffset += w
		s.ch = r
	} else {
		s.offset = len(s.src)
		if s.ch == '\n' {
			s.file.AddLine(s.offset)
		}
		s.ch = eof
	}
}

func (s *Scanner) error(offs int, msg string) {
	if s.errHandler != nil {
		s.errHandler(s.fset.Position(s.file.Pos(offs)), msg)
	}
}

func (s *Scanner) errorf(offs int, format string, args ...any) {
	s.error(offs, fmt.Sprintf(format, args...))
}

// scanComment scans a TeX comment (% comment)
func (s *Scanner) scanComment() string {
	offs := s.offset

	// Scan to the end of the line or file
	for s.ch != '\n' && s.ch != eof {
		s.next()
	}

	return string(s.src[offs:s.offset])
}

// scanCommand scans a TeX command sequence (\command)
func (s *Scanner) scanCommand() (token.Token, string) {
	// Save the starting position of the command (just after the \)
	offs := s.offset

	// Scan the command name
	for isCommandChar(s.ch) {
		s.next()
	}

	// Extract the command name from source (without the \)
	cmdName := string(s.src[offs:s.offset])

	// Look up keyword or return command token (fallback if no keyword)
	return token.LookupKeyword(cmdName), cmdName
}

// scanWord scans a word (sequence of letters)
func (s *Scanner) scanWord() string {
	offs := s.offset

	// Scan the word
	for isLetter(s.ch) {
		s.next()
	}

	// Extract the word from source
	return string(s.src[offs:s.offset])
}

// scanNumber scans a number (integer only)
func (s *Scanner) scanNumber() string {
	offs := s.offset - 1 // -1 to include the first digit

	// Scan the integer part
	for isDigit(s.ch) {
		s.next()
	}

	// Extract the number from source
	return string(s.src[offs:s.offset])
}

// skipWhitespace skips whitespace characters
func (s *Scanner) skipWhitespace() bool {
	skipped := false
	for isSpaceChar(s.ch) {
		s.next()
		skipped = true
	}
	return skipped
}

// Scan scans the next token and returns its position, token type, and literal string
func (s *Scanner) Scan() (pos token.Pos, tok token.Token, lit string) {
	// Skip whitespace
	s.skipWhitespace()

	pos = s.file.Pos(s.offset)

	// Determine token based on the current character
	switch ch := s.ch; {
	case isLetter(ch):
		tok = token.WORD
		lit = s.scanWord()

	case ch == '\\':
		s.next() // consume the \

		switch ch = s.ch; {
		case isCommandChar(ch):
			tok, lit = s.scanCommand()

		case token.IsSymbol(ch):
			s.next()
			tok = token.WORD
			lit = string(ch)

		case isSpaceChar(ch) || ch == '\n':
			// Special case: if \ is followed by a space or newline,
			// it's a special "control space" or escaped newline
			cmd := ""
			if ch == '\n' {
				cmd = "newline"
			} else {
				cmd = "space"
			}
			s.next() // consume the space/newline
			tok, lit = token.COMMAND, cmd

		default:
			s.next()
			tok, lit = token.ILLEGAL, string(ch)
		}

	case ch == '\n':
		s.next()
		tok = token.NEWLINE
		lit = "\n"

	case isDigit(ch):
		s.next()
		tok = token.NUMBER
		lit = s.scanNumber()

	case ch == '%':
		// Comment
		tok = token.COMMENT
		lit = s.scanComment()

	case token.IsSymbol(ch):
		s.next()
		tok = token.LookupSymbol(ch)
		lit = string(ch)

	case ch == eof:
		tok = token.EOF
		lit = "EOF"

	default:
		// Anything else is treated as illegal
		s.next() // consume the character
		tok = token.ILLEGAL
		lit = string(ch)
		s.error(s.offset, fmt.Sprintf("illegal character %#U", ch))
	}

	return
}
