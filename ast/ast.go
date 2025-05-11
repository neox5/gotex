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
	Body     []Node
	Pos_     token.Pos
	End_     token.Pos
}

func (f *File) Pos() token.Pos { return f.Pos_ }
func (f *File) End() token.Pos { return f.End_ }

// ----------------------------------------------------------------------------

type Comment struct {
	Lit        string
	Pos_, End_ token.Pos
}

func (c *Comment) Pos() token.Pos { return c.Pos_ }
func (c *Comment) End() token.Pos { return c.End_ }

// ----------------------------------------------------------------------------

// Word node (e.g., "hello", "world")
type Word struct {
	Lit        string
	Pos_, End_ token.Pos
}

func (w *Word) Pos() token.Pos { return w.Pos_ }
func (w *Word) End() token.Pos { return w.End_ }

// Newline node (single line break, preserved syntactically)
type Newline struct {
	Pos_, End_ token.Pos
}

func (n *Newline) Pos() token.Pos { return n.Pos_ }
func (n *Newline) End() token.Pos { return n.End_ }

// LineBreak (explicit `\\`, `\newline`, etc.)
type LineBreak struct {
	Kind       string
	Pos_, End_ token.Pos
}

func (b *LineBreak) Pos() token.Pos { return b.Pos_ }
func (b *LineBreak) End() token.Pos { return b.End_ }

type TextNode interface {
	Node
	textNode()
}

func (w *Word) textNode()      {}
func (n *Newline) textNode()   {}
func (b *LineBreak) textNode() {}

type TextBlock struct {
	Content    []TextNode
	Pos_, End_ token.Pos
}

func (t *TextBlock) Pos() token.Pos { return t.Pos_ }
func (t *TextBlock) End() token.Pos { return t.End_ }
