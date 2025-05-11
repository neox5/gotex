package ast

import "github.com/neox5/gotex/token"

// All node types impilement the Node interface
type Node interface {
	Pos() token.Pos // position of first character belonging to the node
	End() token.Pos // position of first character immediately after the node
}

// ImportSpec represents \import{}, \input{}, or \usemodule{} commands
type ImportSpec struct {
	Token token.Token // token.IMPORT, token.COMMAND, etc.
	Name  string      // Logical name from braces
	Path  string      // Resolved path (set later)
	Pos_  token.Pos
	End_  token.Pos
}

func (s *ImportSpec) Pos() token.Pos { return s.Pos_ }
func (s *ImportSpec) End() token.Pos { return s.End_ }

// File represents a parsed .tex file in ImportsOnly mode
type File struct {
	Filename string
	Imports  []*ImportSpec
	Pos_     token.Pos
	End_     token.Pos
}

func (f *File) Pos() token.Pos { return f.Pos_ }
func (f *File) End() token.Pos { return f.End_ }
