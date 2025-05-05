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

	// Punctuation
	PERIOD    // .
	COMMA     // ,
	SEMICOLON // ;
	EQUALS    // =
	COLON     // :

	// Keywords begin
	keywords_beg
	KW_IMPORT // \import
	KW_BEGIN  // \begin environment
	KW_END    // \end environment
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

	PERIOD:    ".",
	COMMA:     ",",
	SEMICOLON: ";",
	EQUALS:    "=",
	COLON:     ":",

	KW_IMPORT: "import",
	KW_BEGIN:  "begin",
	KW_END:    "end",
}

// String returns the string representation of the token
func (t Token) String() string {
	if t < 0 || int(t) >= len(tokens) {
		return "token(" + string(rune(t)) + ")"
	}
	return tokens[t]
}

// Add a map to store keywords
var keywords map[string]Token

func init() {
	keywords = make(map[string]Token, keywords_end-(keywords_beg+1))
	for i := keywords_beg + 1; i < keywords_end; i++ {
		// Store without the backslash already
		keywords[tokens[i]] = i
	}
}

// Lookup maps a command name to its keyword token or COMMAND (if not a keyword)
func Lookup(commandName string) Token {
	if tok, isKeyword := keywords[commandName]; isKeyword {
		return tok
	}
	return COMMAND
}
