package ast_test

import (
	"testing"

	"github.com/kamilturek/monke/ast"
	"github.com/kamilturek/monke/token"
)

func TestString(t *testing.T) {
	t.Parallel()

	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{
					Type:    token.LET,
					Literal: "let",
				},
				Name: &ast.Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "myVar",
					},
					Value: "myVar",
				},
				Value: &ast.Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "anotherVar",
					},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("wrong program.String(). got=%q", program.String())
	}
}
