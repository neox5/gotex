package token

import (
	"fmt"
	"unicode"
)

type Token int

const (
	ILLEGAL Token = iota
	EOF
	COMMENT // % This is a comment

	// Commands and structure
	COMMAND  // \documentclass, \begin, \end, etc.
	LBRACE   // { (opening group, command argument)
	RBRACE   // } (closing group, end of argument)
	LBRACKET // [ (for optional arguments)
	RBRACKET // ] (closing optional arguments)

	// Content
	WORD       // Alphabetic sequences (e.g., "Hello", "article")
	NUMBER     // Numeric sequences (e.g., "123", "42")
	WHITESPACE // Space (0x20), tab (0x09)
	NEWLINE    // Line breaks (LF: 0x0A for Unix/Linux, CRLF: 0x0D0A for Windows, CR: 0x0D for classic Mac)

	// Symbols (all punctuation and special symbols)
	symbol_beg

	// Punctuation
	PERIOD    // .
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :

	// Operators
	EQUALS    // =
	LESS      //
	GREATER   // >
	BACKSLASH // literal \
	SLASH     // /
	ASTERISK  // *
	BANG      // !

	// TeX special characters
	AMPERSAND  // &
	DOLLAR     // $
	PERCENT    // %
	HASH       // #
	CARET      // ^
	UNDERSCORE // _
	TILDE      // ~
	PIPE       // |
	AT         // @

	symbol_end

	// Keywords (used for syntax-only primitives)
	keywords_beg
	IMPORT // \import
	ENV    // \begin environment
	ENVEND // \end environment
	keywords_end
)

// tokens maps each token to its string representation for debugging
var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	COMMAND:  "COMMAND",
	LBRACE:   "LBRACE",
	RBRACE:   "RBRACE",
	LBRACKET: "LBRACKET",
	RBRACKET: "RBRACKET",

	WORD:       "WORD",
	NUMBER:     "NUMBER",
	WHITESPACE: "WHITESPACE",
	NEWLINE:    "NEWLINE",

	// Punctuation
	PERIOD:    ".",
	COMMA:     ",",
	SEMICOLON: ";",
	COLON:     ":",

	// Operators
	EQUALS:    "=",
	LESS:      "<",
	GREATER:   ">",
	BACKSLASH: "\\",
	SLASH:     "/",
	ASTERISK:  "*",
	BANG:      "!",

	// TeX special characters
	AMPERSAND:  "&",
	DOLLAR:     "$",
	PERCENT:    "%",
	HASH:       "#",
	CARET:      "^",
	UNDERSCORE: "_",
	TILDE:      "~",
	PIPE:       "|",
	AT:         "@",

	IMPORT: "import",
	ENV:    "begin",
	ENVEND: "end",
}

// Map of symbol runes to their token values
var symbolMap = map[rune]Token{
	// Punctuation
	'.': PERIOD,
	',': COMMA,
	';': SEMICOLON,
	':': COLON,

	// Operators
	'=':  EQUALS,
	'<':  LESS,
	'>':  GREATER,
	'\\': BACKSLASH,
	'/':  SLASH,
	'*':  ASTERISK,
	'!':  BANG,

	// TeX special characters
	'&': AMPERSAND,
	'$': DOLLAR,
	'%': PERCENT,
	'#': HASH,
	'^': CARET,
	'_': UNDERSCORE,
	'~': TILDE,
	'|': PIPE,
	'@': AT,
}

// Map to store keywords
var keywords map[string]Token

func init() {
	// Initialize keywords map
	keywords = make(map[string]Token, keywords_end-(keywords_beg+1))
	for i := keywords_beg + 1; i < keywords_end; i++ {
		// Store without the backslash already
		keywords[tokens[i]] = i
	}

	// Validate the symbolMap
	for r, tok := range symbolMap {
		if tok <= symbol_beg || tok >= symbol_end {
			panic(fmt.Sprintf("token %v for rune %q is not in symbol range", tok, r))
		}
	}
}

// String returns the string representation of the token
func (t Token) String() string {
	if t < 0 || int(t) >= len(tokens) {
		return "token(" + string(rune(t)) + ")"
	}
	return tokens[t]
}

// LookupCommand checks if a command name is a keyword and returns the appropriate token.
// Returns COMMAND for regular commands, or a specific keyword token if it's a keyword.
func LookupCommand(name string) Token {
	if tok, isKeyword := keywords[name]; isKeyword {
		return tok
	}
	return COMMAND
}

// LookupSymbol returns the token for a symbol rune or ILLEGAL if not a symbol.
func LookupSymbol(ch rune) Token {
	if tok, ok := symbolMap[ch]; ok {
		return tok
	}
	return ILLEGAL
}

// IsSymbol reports whether a rune represents a TeX symbol.
func IsSymbol(ch rune) bool {
	_, ok := symbolMap[ch]
	return ok
}

// IsKeyword reports whether tok is a keyword token.
func IsKeyword(tok Token) bool {
	return tok > keywords_beg && tok < keywords_end
}

// IsCommandChar reports whether ch can be part of a TeX command name.
func IsCommandChar(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '@'
}

// IsSpaceChar reports whether ch is a space character.
func IsSpaceChar(ch rune) bool {
	return ch == ' ' || ch == '\t'
}

// IsDigit reports whether ch is a digit.
func IsDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}
