package parser

import (
	"errors"

	"github.com/neox5/gotex/ast"
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
	p := newParser(fset, file, src)

	switch {
	case mode&ImportsOnly != 0:
		return p.parseImportsOnly(), nil
	case mode&ParseFull != 0:
		return p.parseFull(), nil
	default:
		return nil, errors.New("unsupported parse mode")
	}
}
