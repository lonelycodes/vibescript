package parser

import (
	"testing"

	"github.com/lonelycodes/vibescript/ast"
	"github.com/lonelycodes/vibescript/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let blibla = 5018490;
`

	l := lexer.New("let.vibe", input)
	p := New(l)

	program := p.ParseProgram()

	if program == nil {
		t.Fatalf("Program was parsed as nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Expected 3 statements, got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"blibla"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("Expected token literal 'let', got %q", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("Expected a let statement, got %T", stmt)
	}

	if letStmt.Name.Value != name {
		t.Errorf("Expected let stmt name to be '%s', got %s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("Expected token literal to be %s, got %s", name, letStmt.Name.TokenLiteral())
	}

	return true
}
