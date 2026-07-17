package lexer

import (
	"github.com/lonelycodes/vibescript/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := "=+(){},;"
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New("token.vibe", input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type is wrong. expected=%q, actual=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal is wrong. expected=%q, actual=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}

}
