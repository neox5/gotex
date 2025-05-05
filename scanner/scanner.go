package scanner

import (
	"fmt"
	"io"
	"unicode"
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

// next reads the next Unicode character into s.ch and updates positioning
func (s *Scanner) next() {
	if s.rdOffset >= len(s.src) {
		s.offset = s.rdOffset
		s.ch = -1 // EOF marker
		return
	}

	s.offset = s.rdOffset
	r, width := rune(s.src[s.rdOffset]), 1
	if r >= utf8.RuneSelf {
		// Not ASCII - decode the UTF-8 rune
		r, width = utf8.DecodeRune(s.src[s.rdOffset:])
		if r == utf8.RuneError && width == 1 {
			s.error("invalid UTF-8 encoding")
		}
	}
	s.rdOffset += width
	s.ch = r

	// Update line information if we encounter a newline
	if r == '\n' {
		s.file.AddLine(s.offset)
	}
}

// peek returns the byte following the most recently read character without advancing
func (s *Scanner) peek() byte {
	if s.rdOffset >= len(s.src) {
		return 0
	}
	return s.src[s.rdOffset]
}

// error reports an error at the current position
func (s *Scanner) error(msg string) {
	if s.errHandler != nil {
		s.errHandler(s.fset.Position(s.file.Pos(s.offset)), msg)
	}
}

// scanComment scans a TeX comment (% comment)
func (s *Scanner) scanComment() string {
	// Save the starting position of the comment
	offs := s.offset - 1 // -1 to include the % character

	// Scan to the end of the line
	for s.ch != '\n' && s.ch >= 0 {
		s.next()
	}

	// Extract the comment from source
	comment := string(s.src[offs:s.offset])

	return comment
}

// scanCommand scans a TeX command sequence (\command)
func (s *Scanner) scanCommand() (token.Token, string) {
	// Save the starting position of the command (just after the \)
	offs := s.offset

	// Special case: if \ is followed by a space or newline,
	// it's a special "control space" or escaped newline
	if token.IsSpaceChar(s.ch) || s.ch == '\n' {
		cmdName := ""
		if s.ch == '\n' {
			cmdName = "newline"
		} else {
			cmdName = "space"
		}
		s.next() // consume the space/newline
		return token.COMMAND, cmdName
	}

	// Scan the command name
	for token.IsCommandChar(s.ch) {
		s.next()
	}

	// Extract the command name from source (without the \)
	cmdName := string(s.src[offs:s.offset])

	// Look up the command name to determine if it's a keyword
	return token.LookupCommand(cmdName), cmdName
}

// scanWord scans a word (sequence of letters)
func (s *Scanner) scanWord() string {
	offs := s.offset - 1 // -1 to include the first letter

	// Scan the word
	for unicode.IsLetter(s.ch) || s.ch == '-' || s.ch == '\'' {
		s.next()
	}

	// Extract the word from source
	return string(s.src[offs:s.offset])
}

// scanNumber scans a number (integer only)
func (s *Scanner) scanNumber() string {
	offs := s.offset - 1 // -1 to include the first digit

	// Scan the integer part
	for token.IsDigit(s.ch) {
		s.next()
	}

	// Extract the number from source
	return string(s.src[offs:s.offset])
}

// skipWhitespace skips whitespace characters
func (s *Scanner) skipWhitespace() bool {
	skipped := false
	for token.IsSpaceChar(s.ch) {
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
	case ch == -1:
		// End of file
		tok = token.EOF
		lit = "EOF"

	case ch == '%':
		// Comment
		s.next() // consume the %
		tok = token.COMMENT
		lit = s.scanComment()

	case ch == '\\':
		s.next() // consume the \

		switch {
		case s.ch == '\\':
			s.next()
			tok = token.BACKSLASH
			lit = "\\"

		case token.IsSymbol(s.ch):
			// Escaped symbol
			sym := s.ch
			s.next()
			tok = token.WORD
			lit = string(sym)

		default:
			tok, lit = s.scanCommand()
		}

	case token.IsSymbol(ch):
		// Symbol token
		s.next() // consume the symbol
		tok = token.LookupSymbol(ch)
		lit = string(ch)

	case token.IsDigit(ch):
		// Number
		s.next()
		tok = token.NUMBER
		lit = s.scanNumber()

	case unicode.IsLetter(ch):
		// Word
		s.next()
		tok = token.WORD
		lit = s.scanWord()

	default:
		// Anything else is treated as illegal
		s.next() // consume the character
		tok = token.ILLEGAL
		lit = string(ch)
		s.error(fmt.Sprintf("illegal character %#U", ch))
	}

	return
}
