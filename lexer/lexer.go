package lexer

import "github.com/lonelycodes/vibescript/token"

type Lexer struct {
	input        string
	fileName     string
	position     int
	readPosition int
	ch           byte
	line         int
	col          int
}

func New(fileName string, input string) *Lexer {
	l := &Lexer{fileName: fileName, input: input, line: 1, col: 0}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	tok.Line = l.line
	tok.Col = l.col
	switch l.ch {
	case '=':
		tok = l.newToken(token.ASSIGN, string(l.ch))
	case ';':
		tok = l.newToken(token.SEMICOLON, string(l.ch))
	case '(':
		tok = l.newToken(token.LPAREN, string(l.ch))
	case ')':
		tok = l.newToken(token.RPAREN, string(l.ch))
	case ',':
		tok = l.newToken(token.COMMA, string(l.ch))
	case '+':
		tok = l.newToken(token.PLUS, string(l.ch))
	case '{':
		tok = l.newToken(token.LBRACE, string(l.ch))
	case '}':
		tok = l.newToken(token.RBRACE, string(l.ch))
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Line: l.line, Col: l.col}
	}
	l.readChar()
	return tok
}

func (l *Lexer) newToken(t token.TokenType, literal string) token.Token {
	return token.Token{Type: t, Literal: literal, Line: l.line, Col: l.col}
}

func (l *Lexer) readChar() {
	if l.ch == '\n' {
		l.line += 1
		l.col = 1
	}
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
	l.col += 1
}
