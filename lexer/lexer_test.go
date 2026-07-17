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
		expectedLine    int
		expectedCol     int
	}{
		{token.ASSIGN, "=", 1, 1},
		{token.PLUS, "+", 1, 2},
		{token.LPAREN, "(", 1, 3},
		{token.RPAREN, ")", 1, 4},
		{token.LBRACE, "{", 1, 5},
		{token.RBRACE, "}", 1, 6},
		{token.COMMA, ",", 1, 7},
		{token.SEMICOLON, ";", 1, 8},
		{token.EOF, "", 1, 9},
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
		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line is wrong. expected=%d, actual=%d", i, tt.expectedLine, tok.Line)
		}
		if tok.Col != tt.expectedCol {
			t.Fatalf("tests[%d] - column is wrong. expected=%d, actual=%d", i, tt.expectedCol, tok.Col)
		}
	}

}
