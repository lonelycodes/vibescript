package lexer

import (
	"strings"

	"github.com/lonelycodes/vibescript/token"
)

type Lexer struct {
	input        string
	FileName     string
	position     int
	readPosition int
	ch           byte
	line         int
	col          int
	ctxState     int  // 0=idle, 1=saw CTX, 2=saw CTX STRING
	ctxMode      bool // TRUE if we're inside a ctx { ... } body
}

func New(fileName string, input string) *Lexer {
	l := &Lexer{FileName: fileName, input: input, line: 1, col: 0}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	if l.ctxMode {
		return l.nextCtxToken()
	}

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
			tok = l.newToken(token.EQ, "==")
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
	case '"':
		tok = l.readString(tok)
		l.updateCtxState(tok.Type)
		return tok
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Line: l.line, Col: l.col}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			l.updateCtxState(tok.Type)
			return tok
		} else if isDigit(l.ch) {
			tok.Type, tok.Literal = l.readNumber()
			l.updateCtxState(tok.Type)
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL, string(l.ch))
		}
	}
	l.readChar()
	l.updateCtxState(tok.Type)
	return tok
}

func (l *Lexer) nextCtxToken() token.Token {
	l.skipCtxWhiteSpace()
	var tok token.Token
	tok.Line = l.line
	tok.Col = l.col

	if l.ch == 0 {
		// unterminated CTX block
		tok.Type = token.EOF
		tok.Literal = ""
		return tok
	}
	if l.ch == '}' {
		// end of the CTX block
		tok.Type = token.RBRACE
		tok.Literal = "}"
		l.ctxMode = false
		l.readChar()
		return tok
	}

	position := l.position
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	tok.Type = token.CTX_LINE
	tok.Literal = strings.TrimRight(l.input[position:l.position], " \t")
	return tok
}

func (l *Lexer) skipCtxWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		if l.ch == '\n' {
			l.line += 1
			l.col = 0
		}
		l.readChar()
	}
}

func (l *Lexer) updateCtxState(t token.TokenType) {
	switch {
	case t == token.CTX:
		l.ctxState = 1
	case l.ctxState == 1 && t == token.STRING:
		l.ctxState = 2
	case l.ctxState == 2 && t == token.LBRACE:
		l.ctxMode = true
		l.ctxState = 0
	default:
		l.ctxState = 0
	}
}

func (l *Lexer) readString(tok token.Token) token.Token {
	position := l.position
	l.readChar() // read the opening '"'

	for l.ch != '"' {
		if l.ch == 0 || l.ch == '\n' {
			tok.Type = token.ILLEGAL
			tok.Literal = l.input[position:l.position]
			return tok
		}
		if l.ch == '\\' {
			l.readChar()
			if l.ch == 0 {
				continue
			}
		}
		l.readChar()
	}

	tok.Type = token.TokenType(token.STRING)
	tok.Literal = l.input[position+1 : l.position]
	l.readChar() // read the closing '"'

	return tok
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) readNumber() (token.TokenType, string) {
	position := l.position
	tokType := token.TokenType(token.INT)

	l.readDigits()

	if l.ch == '.' && isDigit(l.peekChar()) {
		tokType = token.TokenType(token.FLOAT)
		l.readChar() // read the '.'
		l.readDigits()
	}

	if l.ch == 'e' && isDigit(l.peekChar()) {
		tokType = token.TokenType(token.FLOAT)
		l.readChar() // read the 'e'
		l.readDigits()
	}

	literal := strings.ReplaceAll(l.input[position:l.position], "_", "")
	return tokType, literal
}

func (l *Lexer) readDigits() {
	for isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
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
	for {
		switch l.ch {
		case ' ', '\t', '\r':
			l.readChar()
		case '\n':
			l.line += 1
			l.col = 0
			l.readChar()
		case '#':
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
		default:
			return
		}
	}
}
