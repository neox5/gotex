package token

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

	// Keywords begin
	keywords_beg
	KW_IMPORT // \import
	KW_BEGIN  // \begin environment
	KW_END    // \end environment
	keywords_end
)

// tokens maps each token to its string representation for debugging
var tokens = [...]string{
	ILLEGAL:    "ILLEGAL",
	EOF:        "EOF",
	COMMENT:    "COMMENT",
	COMMAND:    "COMMAND",
	LBRACE:     "LBRACE",
	RBRACE:     "RBRACE",
	LBRACKET:   "LBRACKET",
	RBRACKET:   "RBRACKET",
	WORD:       "WORD",
	NUMBER:     "NUMBER",
	WHITESPACE: "WHITESPACE",
	NEWLINE:    "NEWLINE",

	KW_IMPORT: "\\import",
	KW_BEGIN:  "\\begin",
	KW_END:    "\\end",
}

// String returns the string representation of the token
func (t Token) String() string {
	if t < 0 || int(t) >= len(tokens) {
		return "token(" + string(rune(t)) + ")"
	}
	return tokens[t]
}
