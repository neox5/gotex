package ast

// Visitor is used to traverse an AST.
type Visitor interface {
	Visit(Node) Visitor
}

// Walk walks the AST starting from node, calling v.Visit for each node.
func Walk(v Visitor, node Node) {
	if node == nil {
		return
	}
	if v = v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	case *TextBlock:
		for _, c := range n.Content {
			Walk(v, c)
		}
	case *File:
		for _, c := range n.Body {
			Walk(v, c)
		}
	}
}
