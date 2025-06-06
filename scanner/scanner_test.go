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
			t.Errorf("%s - token %d: expected {%s, %q} [%x], got {%s, %q} [%x] at position %s",
				filename, i, exp.tok, exp.lit, []byte(exp.lit),
				tok, lit, []byte(lit), fset.Position(pos))
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
	% This is a second comment
	\section{Hello} % Inline comment
	Normal text`

	expected := []tokenData{
		{token.COMMENT, "% This is a comment"},
		{token.NEWLINE, "\n"},
		{token.COMMENT, "% This is a second comment"},
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

func TestScanKeywords(t *testing.T) {
	src := `\import{chapter1}
	\begin{matrix}
	\end{matrix}`

	expected := []tokenData{
		{token.IMPORT, "import"},
		{token.LBRACE, "{"},
		{token.WORD, "chapter1"},
		{token.RBRACE, "}"},
		{token.NEWLINE, "\n"},
		{token.ENV, "begin"},
		{token.LBRACE, "{"},
		{token.WORD, "matrix"},
		{token.RBRACE, "}"},
		{token.NEWLINE, "\n"},
		{token.ENVEND, "end"},
		{token.LBRACE, "{"},
		{token.WORD, "matrix"},
		{token.RBRACE, "}"},
		{token.EOF, "EOF"},
	}

	runScannerTest(t, src, expected, "comment_test.tex")
}

func TestWordsWithEscapedSymbols(t *testing.T) {
	src := `foo\$bar a\_b hello\@world`
	expected := []tokenData{
		{token.WORD, "foo$bar"},
		{token.WORD, "a_b"},
		{token.WORD, "hello@world"},
		{token.EOF, "EOF"},
	}
	runScannerTest(t, src, expected, "escaped_in_words_test.tex")
}

func TestEscapedNewlineAndSpace(t *testing.T) {
	src := `word1\newline word2\ word3`
	expected := []tokenData{
		{token.WORD, "word1"},
		{token.COMMAND, "linebreak"},
		{token.WORD, "word2"},
		{token.COMMAND, "space"},
		{token.WORD, "word3"},
		{token.EOF, "EOF"},
	}
	runScannerTest(t, src, expected, "escaped_commands_test.tex")
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
		{token.LBRACK, "["},
		{token.WORD, "Short"},
		{token.WORD, "title"},
		{token.RBRACK, "]"},
		{token.LBRACE, "{"},
		{token.WORD, "Long"},
		{token.WORD, "title"},
		{token.RBRACE, "}"},
		{token.NEWLINE, "\n"},
		{token.COMMAND, "includegraphics"},
		{token.LBRACK, "["},
		{token.WORD, "width"},
		{token.EQUALS, "="},
		{token.NUMBER, "5"},
		{token.WORD, "cm"},
		{token.COMMA, ","},
		{token.WORD, "height"},
		{token.EQUALS, "="},
		{token.NUMBER, "3"},
		{token.WORD, "cm"},
		{token.RBRACK, "]"},
		{token.LBRACE, "{"},
		{token.WORD, "image"},
		{token.PERIOD, "."},
		{token.WORD, "png"},
		{token.RBRACE, "}"},
	}

	runScannerTest(t, src, expected, "optional_args_test.tex")
}
