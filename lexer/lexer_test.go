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

func TestTwoCharOperators(t *testing.T) {
	input := `let x = a |> f();
if x != 1 { ret x // 2; }
y <= 3; y >= 4; b | c; -> ? !
`

	tests := []testCase{
		{token.LET, "let", 1, 1},
		{token.IDENT, "x", 1, 5},
		{token.ASSIGN, "=", 1, 7},
		{token.IDENT, "a", 1, 9},
		{token.PIPEOP, "|>", 1, 11},
		{token.IDENT, "f", 1, 14},
		{token.LPAREN, "(", 1, 15},
		{token.RPAREN, ")", 1, 16},
		{token.SEMICOLON, ";", 1, 17},

		{token.IF, "if", 2, 1},
		{token.IDENT, "x", 2, 4},
		{token.NOT_EQ, "!=", 2, 6},
		{token.INT, "1", 2, 9},
		{token.LBRACE, "{", 2, 11},
		{token.RET, "ret", 2, 13},
		{token.IDENT, "x", 2, 17},
		{token.INTDIV, "//", 2, 19},
		{token.INT, "2", 2, 22},
		{token.SEMICOLON, ";", 2, 23},
		{token.RBRACE, "}", 2, 25},

		{token.IDENT, "y", 3, 1},
		{token.LTE, "<=", 3, 3},
		{token.INT, "3", 3, 6},
		{token.SEMICOLON, ";", 3, 7},
		{token.IDENT, "y", 3, 9},
		{token.GTE, ">=", 3, 11},
		{token.INT, "4", 3, 14},
		{token.SEMICOLON, ";", 3, 15},
		{token.IDENT, "b", 3, 17},
		{token.PIPE, "|", 3, 19},
		{token.IDENT, "c", 3, 21},
		{token.SEMICOLON, ";", 3, 22},
		{token.ARROW, "->", 3, 24},
		{token.QUESTION, "?", 3, 27},
		{token.BANG, "!", 3, 29},

		{token.EOF, "", 4, 1},
	}

	assertLexer(t, input, tests)
}

func TestNoWhitespace(t *testing.T) {
	input := "a1;b(2);"
	tests := []testCase{
		{token.IDENT, "a1", 1, 1},
		{token.SEMICOLON, ";", 1, 3},
		{token.IDENT, "b", 1, 4},
		{token.LPAREN, "(", 1, 5},
		{token.INT, "2", 1, 6},
		{token.RPAREN, ")", 1, 7},
		{token.SEMICOLON, ";", 1, 8},
		{token.EOF, "", 1, 9},
	}
	assertLexer(t, input, tests)
}

func TestKeywordsAndIdentifiers(t *testing.T) {
	input := `fn
let
var
if
elif
else
for
in
while
match
ret
use
ctx
true
false
none
and
or
not
try
err
brk
skip
fnord
letter
iffy
retro
intact
matches
truex
_foo
_
`
	tests := []testCase{
		{token.FN, "fn", 1, 1},
		{token.LET, "let", 2, 1},
		{token.VAR, "var", 3, 1},
		{token.IF, "if", 4, 1},
		{token.ELIF, "elif", 5, 1},
		{token.ELSE, "else", 6, 1},
		{token.FOR, "for", 7, 1},
		{token.IN, "in", 8, 1},
		{token.WHILE, "while", 9, 1},
		{token.MATCH, "match", 10, 1},
		{token.RET, "ret", 11, 1},
		{token.USE, "use", 12, 1},
		{token.CTX, "ctx", 13, 1},
		{token.TRUE, "true", 14, 1},
		{token.FALSE, "false", 15, 1},
		{token.NONE, "none", 16, 1},
		{token.AND, "and", 17, 1},
		{token.OR, "or", 18, 1},
		{token.NOT, "not", 19, 1},
		{token.TRY, "try", 20, 1},
		{token.ERR, "err", 21, 1},
		{token.BRK, "brk", 22, 1},
		{token.SKIP, "skip", 23, 1},

		{token.IDENT, "fnord", 24, 1},
		{token.IDENT, "letter", 25, 1},
		{token.IDENT, "iffy", 26, 1},
		{token.IDENT, "retro", 27, 1},
		{token.IDENT, "intact", 28, 1},
		{token.IDENT, "matches", 29, 1},
		{token.IDENT, "truex", 30, 1},

		{token.IDENT, "_foo", 31, 1},
		{token.IDENT, "_", 32, 1},
		{token.EOF, "", 33, 1},
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
