package ast

import "go/token"

// All node types impilement the Node interface
type Node interface {
	Pos() token.Pos // position of first character belonging to the node
	End() token.Pos // position of first character immediately after the node
}
