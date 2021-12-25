package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers & literals
	IDENT = "IDENT"
	INT   = "INT"

	// Operators
	ASSIGN   = "="
	ASTERISK = "*"
	BANG     = "!"
	MINUS    = "-"
	PLUS     = "+"
	SLASH    = "/"

	EQ     = "=="
	LT     = "<"
	GT     = ">"
	NOT_EQ = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	ELSE     = "ELSE"
	FUNCTION = "FUNCTION"
	FALSE    = "FALSE"
	IF       = "IF"
	LET      = "LET"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"else":   ELSE,
	"false":  FALSE,
	"fn":     FUNCTION,
	"if":     IF,
	"let":    LET,
	"return": RETURN,
	"true":   TRUE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
