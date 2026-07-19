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
	l.skipWhiteSpace()

	var tok token.Token
	tok.Line = l.line
	tok.Col = l.col

	switch l.ch {

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
	case '?':
		tok = l.newToken(token.QUESTION, string(l.ch))
	case '*':
		tok = l.newToken(token.STAR, string(l.ch))
	case '%':
		tok = l.newToken(token.PERCENT, string(l.ch))
	case ':':
		tok = l.newToken(token.COLON, string(l.ch))
	case '[':
		tok = l.newToken(token.LBRACKET, string(l.ch))
	case ']':
		tok = l.newToken(token.RBRACKET, string(l.ch))
	case '.':
		tok = l.newToken(token.DOT, string(l.ch))
	case '@':
		tok = l.newToken(token.AT, string(l.ch))
	case '>':
		if l.peekChar() == '=' {
			tok = l.newToken(token.GTE, ">=")
			l.readChar()
		} else {
			tok = l.newToken(token.GT, string(l.ch))
		}
	case '<':
		if l.peekChar() == '=' {
			tok = l.newToken(token.LTE, "<=")
			l.readChar()
		} else {
			tok = l.newToken(token.LT, string(l.ch))
		}
	case '=':
		if l.peekChar() == '=' {
			tok = l.newToken(token.EQ, "//")
			l.readChar()
		} else {
			tok = l.newToken(token.ASSIGN, string(l.ch))
		}
	case '/':
		if l.peekChar() == '/' {
			tok = l.newToken(token.INTDIV, "//")
			l.readChar()
		} else {
			tok = l.newToken(token.SLASH, string(l.ch))
		}
	case '|':
		if l.peekChar() == '>' {
			tok = l.newToken(token.PIPEOP, "|>")
			l.readChar()
		} else {
			tok = l.newToken(token.PIPE, string(l.ch))
		}
	case '!':
		if l.peekChar() == '=' {
			tok = l.newToken(token.NOT_EQ, "!=")
			l.readChar()
		} else {
			tok = l.newToken(token.BANG, string(l.ch))
		}
	case '-':
		if l.peekChar() == '>' {
			tok = l.newToken(token.ARROW, "->")
			l.readChar()
		} else {
			tok = l.newToken(token.MINUS, string(l.ch))
		}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Line: l.line, Col: l.col}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL, string(l.ch))
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) newToken(t token.TokenType, literal string) token.Token {
	return token.Token{Type: t, Literal: literal, Line: l.line, Col: l.col}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
	l.col += 1
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line += 1
			l.col = 0
		}
		l.readChar()
	}
}
