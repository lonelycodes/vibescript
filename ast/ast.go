package ast

import (
	"strings"

	"github.com/lonelycodes/vibescript/token"
)

type Node interface {
	TokenLiteral() string
	String() string
	Pos() (line, col int)
}

type Statement interface {
	Node
	statementNode()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) String() string {
	return i.Value
}
func (i *Identifier) Pos() (int, int) {
	return i.Token.Line, i.Token.Col
}

type TypeAnnotation struct {
	Token    token.Token
	Name     string
	Elem     *TypeAnnotation // recursive for list[T] / map[T]
	Optional bool            // T?
	Fallible bool            // T! (in return types only)
}

func (t *TypeAnnotation) expressionNode() {}
func (t *TypeAnnotation) TokenLiteral() string {
	return t.Token.Literal
}
func (t *TypeAnnotation) Pos() (int, int) {
	return t.Token.Line, t.Token.Col
}
func (t *TypeAnnotation) String() string {
	var out strings.Builder
	out.WriteString(t.Name)
	if t.Elem != nil {
		out.WriteString("[")
		out.WriteString(t.Elem.String())
		out.WriteString("]")
	}
	if t.Optional {
		out.WriteString("?")
	}
	if t.Fallible {
		out.WriteString("!")
	}
	return out.String()
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
	Type  *TypeAnnotation
}

func (ls *LetStatement) statementNode()
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
func (ls *LetStatement) String() string {
	var out strings.Builder
	out.WriteString(ls.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ls.Name.String())
	if ls.Type != nil {
		out.WriteString(": ")
		out.WriteString(ls.Type.String())
	}
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	return out.String()
}
func (ls *LetStatement) Pos() (int, int) {
	return ls.Token.Line, ls.Token.Col
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
