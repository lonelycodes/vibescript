package lexer

import (
	"github.com/lonelycodes/vibescript/token"
	"testing"
)

type testCase struct {
	expectedType    token.TokenType
	expectedLiteral string
	expectedLine    int
	expectedCol     int
}

func TestBasicFunctionality(t *testing.T) {
	input := "=+(){},;"
	tests := []testCase{
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

	assertLexer(t, input, tests)
}

func TestAddFunction(t *testing.T) {
	input := `let five = 5;
let ten = 10;
let add = |x, y| x + y;
let result = add(five, ten);
`
	tests := []testCase{
		{token.LET, "let", 1, 1},
		{token.IDENT, "five", 1, 5},
		{token.ASSIGN, "=", 1, 10},
		{token.INT, "5", 1, 12},
		{token.SEMICOLON, ";", 1, 13},

		{token.LET, "let", 2, 1},
		{token.IDENT, "ten", 2, 5},
		{token.ASSIGN, "=", 2, 9},
		{token.INT, "10", 2, 11},
		{token.SEMICOLON, ";", 2, 13},

		{token.LET, "let", 3, 1},
		{token.IDENT, "add", 3, 5},
		{token.ASSIGN, "=", 3, 9},
		{token.PIPE, "|", 3, 11},
		{token.IDENT, "x", 3, 12},
		{token.COMMA, ",", 3, 13},
		{token.IDENT, "y", 3, 15},
		{token.PIPE, "|", 3, 16},
		{token.IDENT, "x", 3, 18},
		{token.PLUS, "+", 3, 20},
		{token.IDENT, "y", 3, 22},
		{token.SEMICOLON, ";", 3, 23},

		{token.LET, "let", 4, 1},
		{token.IDENT, "result", 4, 5},
		{token.ASSIGN, "=", 4, 12},
		{token.IDENT, "add", 4, 14},
		{token.LPAREN, "(", 4, 17},
		{token.IDENT, "five", 4, 18},
		{token.COMMA, ",", 4, 22},
		{token.IDENT, "ten", 4, 24},
		{token.RPAREN, ")", 4, 27},
		{token.SEMICOLON, ";", 4, 28},

		{token.EOF, "", 5, 1},
	}

	assertLexer(t, input, tests)
}

func assertLexer(t *testing.T, input string, tests []testCase) {
	t.Helper()
	l := New("test.vibe", input)
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
