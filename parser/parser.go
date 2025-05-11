package parser

import (
	"errors"

	"github.com/neox5/gotex/ast"
	"github.com/neox5/gotex/scanner"
	"github.com/neox5/gotex/token"
)

// Mode controls which parsing features are enabled.
type Mode uint

const (
	ImportsOnly Mode = 1 << iota // Only parse \import statements
	ParseFull                    // Future: parse full syntax tree
)

// Parse parses the given source into a syntax tree depending on the mode.
// The source must be valid UTF-8. The caller must provide a token.FileSet and associated token.File.
func Parse(fset *token.FileSet, file *token.File, src []byte, mode Mode) (*ast.File, error) {
	var scan scanner.Scanner
	scan.Init(fset, file, src, nil)

	if mode&ImportsOnly != 0 {
		return parseImportsOnly(file, src, &scan), nil
	}

	return nil, errors.New("ParseFull mode not implemented")
}

// parseImportsOnly is a NOOP stub that returns an empty AST with correct position info.
func parseImportsOnly(file *token.File, src []byte, _ *scanner.Scanner) *ast.File {
	return &ast.File{
		Filename: file.Name(),
		Imports:  []*ast.ImportSpec{},
		Pos_:     file.Pos(0),
		End_:     file.Pos(len(src)),
	}
}
