package parser

import (
	"github.com/lonelycodes/vibescript/ast"
	"github.com/lonelycodes/vibescript/lexer"
	"github.com/lonelycodes/vibescript/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	// reading two tokens so both current and peek are set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}
