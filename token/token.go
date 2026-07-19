package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
}

const (
	// special
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// identifiers & literals
	IDENT  = "IDENT"
	INT    = "INT"
	FLOAT  = "FLOAT"
	STRING = "STRING"

	// operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	STAR     = "*"
	SLASH    = "/"
	INTDIV   = "//"
	PERCENT  = "%"
	EQ       = "=="
	NOT_EQ   = "!="
	LT       = "<"
	LTE      = "<="
	GT       = ">"
	GTE      = ">="
	PIPEOP   = "|>"
	ARROW    = "->"
	QUESTION = "?"
	BANG     = "!"
	PIPE     = "|"
	DOT      = "."

	// delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"
	AT        = "@"

	// keywords
	FN       = "FN"
	LET      = "LET"
	VAR      = "VAR"
	IF       = "IF"
	ELIF     = "ELIF"
	ELSE     = "ELSE"
	FOR      = "FOR"
	IN       = "IN"
	WHILE    = "WHILE"
	MATCH    = "MATCH"
	RET      = "RET"
	USE      = "USE"
	CTX      = "CTX"
	CTX_LINE = "CTX_LINE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	NONE     = "NONE"
	AND      = "AND"
	OR       = "OR"
	NOT      = "NOT"
	TRY      = "TRY"
	ERR      = "ERR"
	BRK      = "BRK"
	SKIP     = "SKIP"
)

var keywords = map[string]TokenType{
	"fn":    FN,
	"let":   LET,
	"var":   VAR,
	"if":    IF,
	"elif":  ELIF,
	"else":  ELSE,
	"for":   FOR,
	"in":    IN,
	"while": WHILE,
	"match": MATCH,
	"ret":   RET,
	"use":   USE,
	"ctx":   CTX,
	"true":  TRUE,
	"false": FALSE,
	"none":  NONE,
	"and":   AND,
	"or":    OR,
	"not":   NOT,
	"try":   TRY,
	"err":   ERR,
	"brk":   BRK,
	"skip":  SKIP,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
