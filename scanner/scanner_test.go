package scanner

import (
	"testing"

	"github.com/neox5/gotex/token"
)

type tokenData struct {
	tok token.Token
	lit string
}

// Helper function to run scanner tests
func runScannerTest(t *testing.T, src string, expected []tokenData, filename string) {
	// Set up the scanner
	fset := token.NewFileSet()
	file := fset.AddFile(filename, fset.Base(), len(src))
	var s Scanner
	s.Init(fset, file, []byte(src), nil)

	// Scan all tokens and compare with expected
	for i, exp := range expected {
		pos, tok, lit := s.Scan()
		if tok != exp.tok || lit != exp.lit {
			t.Errorf("%s - token %d: expected {%s, %q}, got {%s, %q} at position %s",
				filename, i, exp.tok, exp.lit, tok, lit, fset.Position(pos))
		}
	}

	// Check that we don't have more tokens
	pos, tok, lit := s.Scan()
	if tok != token.EOF {
		t.Errorf("%s - expected EOF, got {%s, %q} at position %s",
			filename, tok, lit, fset.Position(pos))
	}
}

func TestScanComments(t *testing.T) {
	src := `% This is a comment
	\section{Hello} % Inline comment
	Normal text`

	expected := []tokenData{
		{token.COMMENT, "% This is a comment"},
		{token.NEWLINE, "\n"},
		{token.COMMAND, "section"},
		{token.LBRACE, "{"},
		{token.WORD, "Hello"},
		{token.RBRACE, "}"},
		{token.COMMENT, "% Inline comment"},
		{token.NEWLINE, "\n"},
		{token.WORD, "Normal"},
		{token.WORD, "text"},
		{token.EOF, "EOF"},
	}

	runScannerTest(t, src, expected, "comment_test.tex")
}

func TestEscapedSymbols(t *testing.T) {
	src := `\\ \ \newline \@\$\%\&\#\_\{\}\~\^`

	expected := []tokenData{
		{token.BACKSLASH, "\\"},
		{token.COMMAND, "space"},
		{token.COMMAND, "newline"},
		{token.WORD, "@"},
		{token.WORD, "$"},
		{token.WORD, "%"},
		{token.WORD, "&"},
		{token.WORD, "#"},
		{token.WORD, "_"},
		{token.WORD, "{"},
		{token.WORD, "}"},
		{token.WORD, "~"},
		{token.WORD, "^"},
		{token.EOF, "EOF"},
	}

	runScannerTest(t, src, expected, "special_commands_test.tex")
}

func TestScanCommands(t *testing.T) {
	src := `\section{Title}
\begin{document}
\end{document}
\normal
\import{path}`

	expected := []tokenData{
		{token.COMMAND, "section"},
		{token.LBRACE, "{"},
		{token.WORD, "Title"},
		{token.RBRACE, "}"},
		{token.NEWLINE, "\n"},
		{token.ENV, "begin"},
		{token.LBRACE, "{"},
		{token.WORD, "document"},
		{token.RBRACE, "}"},
		{token.NEWLINE, "\n"},
		{token.ENVEND, "end"},
		{token.LBRACE, "{"},
		{token.WORD, "document"},
		{token.RBRACE, "}"},
		{token.NEWLINE, "\n"},
		{token.COMMAND, "normal"},
		{token.NEWLINE, "\n"},
		{token.IMPORT, "import"},
		{token.LBRACE, "{"},
		{token.WORD, "path"},
		{token.RBRACE, "}"},
	}

	runScannerTest(t, src, expected, "command_test.tex")
}

func TestScanOptionalArguments(t *testing.T) {
	src := `\section[Short title]{Long title}
\includegraphics[width=5cm, height=3cm]{image.png}`

	expected := []tokenData{
		{token.COMMAND, "section"},
		{token.LBRACKET, "["},
		{token.WORD, "Short"},
		{token.WORD, "title"},
		{token.RBRACKET, "]"},
		{token.LBRACE, "{"},
		{token.WORD, "Long"},
		{token.WORD, "title"},
		{token.RBRACE, "}"},
		{token.NEWLINE, "\n"},
		{token.COMMAND, "includegraphics"},
		{token.LBRACKET, "["},
		{token.WORD, "width"},
		{token.EQUALS, "="},
		{token.NUMBER, "5"},
		{token.WORD, "cm"},
		{token.COMMA, ","},
		{token.WORD, "height"},
		{token.EQUALS, "="},
		{token.NUMBER, "3"},
		{token.WORD, "cm"},
		{token.RBRACKET, "]"},
		{token.LBRACE, "{"},
		{token.WORD, "image"},
		{token.PERIOD, "."},
		{token.WORD, "png"},
		{token.RBRACE, "}"},
	}

	runScannerTest(t, src, expected, "optional_args_test.tex")
}

func TestScanNumbersAndWords(t *testing.T) {
	src := `Simple words and numbers 123 45.67
Words-with-hyphens and apostrophe's are 
treated as one word. Numbers like 3.14159 are parsed.`

	expected := []tokenData{
		{token.WORD, "Simple"},
		{token.WORD, "words"},
		{token.WORD, "and"},
		{token.WORD, "numbers"},
		{token.NUMBER, "123"},
		{token.NUMBER, "45"},
		{token.PERIOD, "."},
		{token.NUMBER, "67"},
		{token.NEWLINE, "\n"},
		{token.WORD, "Words-with-hyphens"},
		{token.WORD, "and"},
		{token.WORD, "apostrophe's"},
		{token.WORD, "are"},
		{token.NEWLINE, "\n"},
		{token.WORD, "treated"},
		{token.WORD, "as"},
		{token.WORD, "one"},
		{token.WORD, "word"},
		{token.PERIOD, "."},
		{token.WORD, "Numbers"},
		{token.WORD, "like"},
		{token.NUMBER, "3"},
		{token.PERIOD, "."},
		{token.NUMBER, "14159"},
		{token.WORD, "are"},
		{token.WORD, "parsed"},
		{token.PERIOD, "."},
	}

	runScannerTest(t, src, expected, "numbers_words_test.tex")
}
