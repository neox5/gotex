package parser

import (
	"os"
	"testing"

	"github.com/neox5/gotex/ast"
	"github.com/neox5/gotex/token"
)

func TestParseWordsOnlyFile(t *testing.T) {
	// Load file content from testdata
	src, err := os.ReadFile("./testdata/words_only.gtex")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	// Prepare file set and token.File
	fset := token.NewFileSet()
	file := fset.AddFile("words_only.gtex", fset.Base(), len(src))

	// Parse in full mode
	astFile, err := Parse(fset, file, src, ParseFull)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if astFile == nil || len(astFile.Body) == 0 {
		t.Fatalf("parsed file is empty or nil")
	}

	// Log output for inspection
	for i, node := range astFile.Body {
		t.Logf("Node %d: %T at %s", i, node, fset.Position(node.Pos()))
	}
}

func TestWordsOnlyCompare(t *testing.T) {
	src, err := os.ReadFile("./testdata/words_only.gtex")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	fset := token.NewFileSet()
	file := fset.AddFile("words_only.gtex", fset.Base(), len(src))
	astFile, err := Parse(fset, file, src, ParseFull)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	expected := &ast.File{
		Filename: "words_only.gtex",
		Body: []ast.Node{
			&ast.Comment{Lit: "% words_only.gtex"},
			&ast.Newline{},
			&ast.Comment{Lit: "% This file demonstrates newline behavior in Gotex"},
			&ast.Newline{},
			&ast.Comment{Lit: "% using plain words and implicit layout — no macros or punctuation."},
			&ast.Newline{},
			&ast.Comment{Lit: "% ------------------------------------------------"},
			&ast.Newline{},
			&ast.Comment{Lit: "% Variant 1: Double newline (blank line)"},
			&ast.Newline{},
			&ast.Comment{Lit: "% This creates a new paragraph in LaTeX/Gotex."},
			&ast.Newline{},
			&ast.Comment{Lit: "% Paragraphs are separated by vertical space."},
			&ast.Newline{},
			&ast.Comment{Lit: "% ------------------------------------------------"},
			&ast.Newline{},

			&ast.TextBlock{Content: []ast.TextNode{
				&ast.Word{Lit: "hello"},
				&ast.Word{Lit: "world"},
				&ast.Newline{},
			}},
			&ast.Newline{},
			&ast.TextBlock{Content: []ast.TextNode{
				&ast.Word{Lit: "this"},
				&ast.Word{Lit: "is"},
				&ast.Word{Lit: "a"},
				&ast.Word{Lit: "new"},
				&ast.Word{Lit: "paragraph"},
				&ast.Word{Lit: "using"},
				&ast.Word{Lit: "a"},
				&ast.Word{Lit: "double"},
				&ast.Word{Lit: "newline"},
				&ast.Newline{},
			}},
			&ast.Newline{},

			&ast.Comment{Lit: "% ------------------------------------------------"},
			&ast.Newline{},
			&ast.Comment{Lit: "% Variant 2: Single newline"},
			&ast.Newline{},
			&ast.Comment{Lit: "% LaTeX and Gotex ignore single newlines and treat them as spaces."},
			&ast.Newline{},
			&ast.Comment{Lit: "% The result is one continuous paragraph."},
			&ast.Newline{},
			&ast.Comment{Lit: "% ------------------------------------------------"},
			&ast.Newline{},

			&ast.TextBlock{Content: []ast.TextNode{
				&ast.Word{Lit: "this"},
				&ast.Word{Lit: "is"},
				&ast.Word{Lit: "a"},
				&ast.Word{Lit: "single"},
				&ast.Word{Lit: "line"},
				&ast.Newline{},
				&ast.Word{Lit: "broken"},
				&ast.Word{Lit: "across"},
				&ast.Word{Lit: "multiple"},
				&ast.Newline{},
				&ast.Word{Lit: "source"},
				&ast.Word{Lit: "lines"},
				&ast.Word{Lit: "but"},
				&ast.Word{Lit: "rendered"},
				&ast.Newline{},
				&ast.Word{Lit: "as"},
				&ast.Word{Lit: "a"},
				&ast.Word{Lit: "single"},
				&ast.Word{Lit: "paragraph"},
				&ast.Newline{},
			}},
			&ast.Newline{},

			&ast.Comment{Lit: "% ------------------------------------------------"},
			&ast.Newline{},
			&ast.Comment{Lit: "% Variant 3: Forced line break using `\\\\`"},
			&ast.Newline{},
			&ast.Comment{Lit: "% This causes an explicit line break within a paragraph."},
			&ast.Newline{},
			&ast.Comment{Lit: "% It is like pressing \"Enter\" but without ending the paragraph."},
			&ast.Newline{},
			&ast.Comment{Lit: "% ------------------------------------------------"},
			&ast.Newline{},

			&ast.TextBlock{Content: []ast.TextNode{
				&ast.Word{Lit: "this"},
				&ast.Word{Lit: "line"},
				&ast.Word{Lit: "ends"},
				&ast.Word{Lit: "here"},
				&ast.LineBreak{Kind: "newline"},
			}},
			&ast.TextBlock{Content: []ast.TextNode{
				&ast.Word{Lit: "and"},
				&ast.Word{Lit: "this"},
				&ast.Word{Lit: "starts"},
				&ast.Word{Lit: "on"},
				&ast.Word{Lit: "the"},
				&ast.Word{Lit: "next"},
				&ast.Word{Lit: "line"},
				&ast.LineBreak{Kind: "newline"},
			}},
			&ast.TextBlock{Content: []ast.TextNode{
				&ast.Word{Lit: "still"},
				&ast.Word{Lit: "within"},
				&ast.Word{Lit: "the"},
				&ast.Word{Lit: "same"},
				&ast.Word{Lit: "paragraph"},
			}},
			&ast.Newline{},

			&ast.Comment{Lit: "% ------------------------------------------------"},
			&ast.Newline{},
			&ast.Comment{Lit: "% Summary:"},
			&ast.Newline{},
			&ast.Comment{Lit: "% - Double newline → new paragraph"},
			&ast.Newline{},
			&ast.Comment{Lit: "% - Single newline → treated as space"},
			&ast.Newline{},
			&ast.Comment{Lit: "% - `\\\\` → line break, same paragraph"},
			&ast.Newline{},
			&ast.Comment{Lit: "% Gotex implicitly wraps all this in:"},
			&ast.Newline{},
			&ast.Comment{Lit: "%"},
			&ast.Newline{},
			&ast.Comment{Lit: "% \\documentclass{article}"},
			&ast.Newline{},
			&ast.Comment{Lit: "% \\begin{document}"},
			&ast.Newline{},
			&ast.Comment{Lit: "% ..."},
			&ast.Newline{},
			&ast.Comment{Lit: "% \\end{document}"},
			&ast.Newline{},
			&ast.Comment{Lit: "% ------------------------------------------------"},
			&ast.Newline{},
		},
	}

	visitor := &ast.CompareVisitor{
		T:               t,
		Expected:        expected,
		SkipLengthCheck: true,
	}
	ast.Walk(visitor, astFile)
	visitor.Finish()
}
