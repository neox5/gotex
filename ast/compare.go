package ast

import (
	"fmt"
	"testing"
)

// CompareVisitor walks an AST and compares it to an expected list of nodes.
type CompareVisitor struct {
	T               *testing.T
	Expected        Node
	SkipLengthCheck bool
	matched         bool // internal flag to track if comparison happened
}

func (v *CompareVisitor) Visit(n Node) Visitor {
	if v.matched {
		// Ignore further nodes after root
		return nil
	}
	v.matched = true

	if !v.compareNodes(n, v.Expected) {
		v.T.Errorf("AST mismatch:\n  got:  %s\n  want: %s", shortNode(n), shortNode(v.Expected))
	}

	return v // still walk children to trigger deeper comparison
}

func (v *CompareVisitor) Finish() {
	if !v.matched {
		v.T.Error("Visitor did not match any root node")
	}
}

// compareNodes compares two AST nodes for structural equality.
func (v *CompareVisitor) compareNodes(a, b Node) bool {
	switch x := a.(type) {
	case *File:
		y, ok := b.(*File)
		if !ok {
			v.T.Errorf("expected File node, got %T", b)
			return false
		}
		if !v.SkipLengthCheck && len(x.Body) != len(y.Body) {
			v.T.Errorf("File.Body length mismatch: got %d nodes, want %d nodes", len(x.Body), len(y.Body))
			return false
		}
		minLen := min(len(x.Body), len(y.Body))
		for i := 0; i < minLen; i++ {
			if !v.compareNodes(x.Body[i], y.Body[i]) {
				v.T.Errorf("File.Body[%d] mismatch:\n  got:  %s\n  want: %s",
					i, shortNode(x.Body[i]), shortNode(y.Body[i]))
				return false
			}
		}
		return true

	case *TextBlock:
		y, ok := b.(*TextBlock)
		if !ok {
			v.T.Errorf("expected TextBlock, got %T", b)
			return false
		}
		if !v.SkipLengthCheck && len(x.Content) != len(y.Content) {
			v.T.Errorf("TextBlock.Content length mismatch: got %d nodes, want %d nodes", len(x.Content), len(y.Content))
			return false
		}
		minLen := min(len(x.Content), len(y.Content))
		for i := 0; i < minLen; i++ {
			if !v.compareNodes(x.Content[i], y.Content[i]) {
				v.T.Errorf("TextBlock.Content[%d] mismatch:\n  got:  %s\n  want: %s",
					i, shortNode(x.Content[i]), shortNode(y.Content[i]))
				return false
			}
		}
		return true

	case *Word:
		y, ok := b.(*Word)
		if !ok || x.Lit != y.Lit {
			v.T.Errorf("Word mismatch: got %q, want %q", x.Lit, y.Lit)
			return false
		}
		return true

	case *Newline:
		_, ok := b.(*Newline)
		if !ok {
			v.T.Errorf("expected Newline, got %T", b)
			return false
		}
		return true

	case *LineBreak:
		y, ok := b.(*LineBreak)
		if !ok || x.Kind != y.Kind {
			v.T.Errorf("LineBreak mismatch: got %q, want %q", x.Kind, y.Kind)
			return false
		}
		return true

	case *Comment:
		y, ok := b.(*Comment)
		if !ok || x.Lit != y.Lit {
			v.T.Errorf("Comment mismatch: got %q, want %q", x.Lit, y.Lit)
			return false
		}
		return true

	default:
		v.T.Errorf("unexpected node type: %T", a)
		return false
	}
}

// shortNode returns a brief string representation of a node.
func shortNode(n Node) string {
	switch x := n.(type) {
	case *File:
		snippets := make([]string, 0, 2)
		for i := 0; i < min(2, len(x.Body)); i++ {
			snippets = append(snippets, shortNode(x.Body[i]))
		}
		more := ""
		if len(x.Body) > 2 {
			more = ", ..."
		}
		return fmt.Sprintf("File[%d nodes: %s%s]", len(x.Body), joinSnippets(snippets), more)

	case *TextBlock:
		snippets := make([]string, 0, 2)
		for i := 0; i < min(2, len(x.Content)); i++ {
			snippets = append(snippets, shortNode(x.Content[i]))
		}
		more := ""
		if len(x.Content) > 2 {
			more = ", ..."
		}
		return fmt.Sprintf("TextBlock[%d: %s%s]", len(x.Content), joinSnippets(snippets), more)

	case *Word:
		return fmt.Sprintf("Word(%q)", x.Lit)
	case *Newline:
		return "Newline"
	case *LineBreak:
		return fmt.Sprintf("LineBreak(%q)", x.Kind)
	case *Comment:
		return fmt.Sprintf("Comment(%q)", x.Lit)
	default:
		return fmt.Sprintf("%T", x)
	}
}

func joinSnippets(parts []string) string {
	return fmt.Sprintf("%s", parts[0]) + func() string {
		if len(parts) == 2 {
			return ", " + parts[1]
		}
		return ""
	}()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
