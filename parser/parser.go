package parser

import (
	"github.com/neox5/gotex/ast"
	"github.com/neox5/gotex/scanner"
	"github.com/neox5/gotex/token"
)

type parser struct {
	s    *scanner.Scanner
	fset *token.FileSet
	file *token.File

	tok token.Token
	lit string
	pos token.Pos
}

func newParser(fset *token.FileSet, file *token.File, src []byte) *parser {
	var scan scanner.Scanner
	scan.Init(fset, file, src, nil)
	p := &parser{
		s:    &scan,
		fset: fset,
		file: file,
	}
	p.next()
	return p
}

func (p *parser) next() {
	p.pos, p.tok, p.lit = p.s.Scan()
}

func (p *parser) parseFull() *ast.File {
	var nodes []ast.Node
	start := p.pos

	for p.tok != token.EOF {
		switch p.tok {
		case token.COMMENT:
			comment := p.parseComment()
			nodes = append(nodes, comment)
			if p.tok == token.NEWLINE {
				nodes = append(nodes, &ast.Newline{
					Pos_: p.pos,
					End_: p.pos + 1,
				})
				p.next()
			}
		case token.NEWLINE:
			text := p.parseText() // treat as part of text
			nodes = append(nodes, text)
		case token.COMMAND:
			if p.lit == "newline" {
				text := p.parseText() // same, groupable
				nodes = append(nodes, text)
			} else {
				// TODO: dispatch to command handling
				p.next()
			}
		case token.WORD:
			text := p.parseText()
			nodes = append(nodes, text)
		default:
			p.next() // skip unknown or unexpected tokens
		}
	}

	return &ast.File{
		Filename: p.file.Name(),
		Imports:  nil, // Not collected in full mode
		Body:     nodes,
		Pos_:     start,
		End_:     p.pos,
	}
}

func (p *parser) parseImportsOnly() *ast.File {
	var imports []*ast.ImportSpec

	for p.tok != token.EOF {
		if p.tok == token.IMPORT {
			imp := p.parseImportSpec()
			if imp != nil {
				imports = append(imports, imp)
			}
		} else {
			p.next() // skip other tokens
		}
	}

	start := token.Pos(p.file.Base())
	end := p.pos

	return &ast.File{
		Filename: p.file.Name(),
		Imports:  imports,
		Pos_:     start,
		End_:     end,
	}
}

func (p *parser) parseImportSpec() *ast.ImportSpec {
	start := p.pos
	cmdTok := p.tok
	p.next() // consume \import

	if p.tok != token.LBRACE {
		// Incomplete import, skip
		return nil
	}
	p.next() // consume {

	var name string
	if p.tok == token.WORD || p.tok == token.COMMAND {
		name = p.lit
		p.next()
		if p.tok == token.NUMBER {
			name += p.lit
			p.next()
		}
	}

	if p.tok != token.RBRACE {
		// Malformed import
		return nil
	}
	end := p.pos
	p.next() // consume }

	return &ast.ImportSpec{
		Token: cmdTok,
		Name:  name,
		Pos_:  start,
		End_:  end + 1,
	}
}

func (p *parser) parseComment() *ast.Comment {
	comment := &ast.Comment{
		Lit:  p.lit,
		Pos_: p.pos,
		End_: p.pos + token.Pos(len(p.lit)),
	}
	p.next()
	return comment
}

func (p *parser) parseText() *ast.TextBlock {
	var content []ast.TextNode
	start := p.pos

loop:
	for p.tok != token.EOF {
		switch p.tok {
		case token.WORD:
			node := &ast.Word{
				Lit:  p.lit,
				Pos_: p.pos,
				End_: p.pos + token.Pos(len(p.lit)),
			}
			content = append(content, node)
			p.next()

		case token.NEWLINE:
			node := &ast.Newline{
				Pos_: p.pos,
				End_: p.pos + 1,
			}
			content = append(content, node)
			p.next()

		case token.COMMAND:
			if p.lit == "newline" {
				node := &ast.LineBreak{
					Kind: "newline",
					Pos_: p.pos,
					End_: p.pos + token.Pos(len(p.lit)),
				}
				content = append(content, node)
				p.next()
			} else {
				break loop // ✅ exits the for-loop
			}

		default:
			break loop // ✅ exits the for-loop on any non-text token
		}
	}

	return &ast.TextBlock{
		Content: content,
		Pos_:    start,
		End_:    p.pos,
	}
}
