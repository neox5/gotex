package token

type Token int

const (
	ILLEGAL Token = iota
	EOF
	COMMENT // % This is a comment

	// Content
	WORD       // Alphabetic sequences (e.g., "Hello", "article")
	NUMBER     // Numeric sequences (e.g., "123", "42")
	WHITESPACE // Space (0x20), tab (0x09)
	NEWLINE    // Line breaks (LF: 0x0A for Unix/Linux, CRLF: 0x0D0A for Windows, CR: 0x0D for classic Mac)

	COMMAND // \documentclass, \begin, \end, etc.

	keywords_beg
	IMPORT // \import
	ENV    // \begin environment
	ENVEND // \end environment
	keywords_end

	symbols_beg
	LBRACE // {
	RBRACE // }
	LBRACK // [
	RBRACK // ]

	PERIOD    // .
	COLON     // :
	COMMA     // ,
	SEMICOLON // ;

	EQUALS    // =
	LESS      // <
	GREATER   // >
	BACKSLASH // \
	SLASH     // /
	ASTERISK  // *
	BANG      // !

	AMPERSAND  // &
	DOLLAR     // $
	PERCENT    // %
	HASH       // #
	CARET      // ^
	DASH       // -
	UNDERSCORE // _
	TILDE      // ~
	PIPE       // |
	AT         // @
	symbols_end
)

// tokens maps each token to its string representation for debugging
var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	WORD:       "WORD",
	NUMBER:     "NUMBER",
	WHITESPACE: "WHITESPACE",
	NEWLINE:    "NEWLINE",

	COMMAND: "COMMAND",

	IMPORT: "import",
	ENV:    "begin",
	ENVEND: "end",

	LBRACE: "{",
	RBRACE: "}",
	LBRACK: "[",
	RBRACK: "]",

	// Punctuation
	PERIOD:    ".",
	COLON:     ":",
	COMMA:     ",",
	SEMICOLON: ";",

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
	DASH:       "-",
	UNDERSCORE: "_",
	TILDE:      "~",
	PIPE:       "|",
	AT:         "@",
}

// Map to store keywords
var (
	keywords map[string]Token
	symbols  map[string]Token
)

func init() {
	// Initialize keywords map
	keywords = make(map[string]Token, keywords_end-(keywords_beg+1))
	for i := keywords_beg + 1; i < keywords_end; i++ {
		keywords[tokens[i]] = i
	}

	// Initialize symbols map
	symbols = make(map[string]Token, symbols_end-(symbols_beg+1))
	for i := symbols_beg + 1; i < symbols_end; i++ {
		symbols[tokens[i]] = i
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
func LookupKeyword(name string) Token {
	if tok, isKeyword := keywords[name]; isKeyword {
		return tok
	}
	return COMMAND
}

// LookupSymbol returns the token for a symbol rune or ILLEGAL if not a symbol.
func LookupSymbol(ch rune) Token {
	if tok, ok := symbols[string(ch)]; ok {
		return tok
	}
	return ILLEGAL
}

// IsSymbol reports whether a rune represents a symbol.
func IsSymbol(ch rune) bool {
	_, ok := symbols[string(ch)]
	return ok
}

// IsKeyword reports whether tok is a keyword token.
func IsKeyword(tok Token) bool {
	return tok > keywords_beg && tok < keywords_end
}
