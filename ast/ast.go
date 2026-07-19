package ast

import "strings"

type Node interface {
	TokenLiteral() string
	String() string
	Pos() (line, col int)
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) Pos() (int, int) {
	if len(p.Statements) > 0 {
		return p.Statements[0].Pos()
	}
	return 1, 1
}

func (p *Program) String() string {
	var out strings.Builder
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
