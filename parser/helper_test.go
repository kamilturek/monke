package parser_test

import (
	"fmt"
	"testing"

	"github.com/kamilturek/monke/ast"
	"github.com/kamilturek/monke/parser"
)

func testInfixExpression(
	t *testing.T, exp ast.Expression, expectedLeft interface{}, expectedOperator string, expectedRight interface{},
) bool {
	t.Helper()

	infixExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("wrong expression type. expected=*ast.InfixExpression. got=%T", exp)
	}

	if !testLiteralExpression(t, infixExp.Left, expectedLeft) {
		return false
	}

	if infixExp.Operator != expectedOperator {
		t.Errorf("wrong operator. expected=%s, got=%s", infixExp.Operator, expectedOperator)
		return false
	}

	if !testLiteralExpression(t, infixExp.Right, expectedRight) {
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	t.Helper()

	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	t.Fatalf("unsupported literal type. got=%T", expected)

	return false
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {
	t.Helper()

	literal, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("stmt.Expression not *ast.IntegerLiteral. got=%T", exp)
		return false
	}

	if literal.Value != value {
		t.Errorf("integerLiteral.Value not %d. got=%d", value, literal.Value)
		return false
	}

	if literal.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integerLiteral.TokenLiteral not %s. got=%s", fmt.Sprintf("%d", value), literal.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	t.Helper()

	boolean, ok := exp.(*ast.Boolean)
	if !ok {
		t.Fatalf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if boolean.Value != value {
		t.Fatalf("boolean.Value not %t. got=%t", value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Fatalf("boolean.TokenLiteral not %t. got=%s", value, boolean.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	t.Helper()

	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Fatalf("wrong type. expected=*ast.Identifier, got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Fatalf("wrong value. expected=%s, got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Fatalf("wrong token literal. expected=%s, got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	t.Helper()

	if s.TokenLiteral() != "let" {
		t.Errorf("wrong token literal. expected=%s got=%s", "let", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("wrong type. expected=*ast.LetStatement, got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("wrong value. expected=%s, got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("wrong token literal. expected=%s, got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func testReturnStatement(t *testing.T, s ast.Statement) bool {
	t.Helper()

	if s.TokenLiteral() != "return" {
		t.Errorf("wrong token literal. expected=%s got=%s", "return", s.TokenLiteral())
		return false
	}

	_, ok := s.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("wrong type. expected=*ast.ReturnStatement, got=%T", s)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	t.Helper()

	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}
