package parser_test

import (
	"testing"

	"github.com/kamilturek/monke/ast"
	"github.com/kamilturek/monke/lexer"
	"github.com/kamilturek/monke/parser"
)

func TestLetStatements(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			l := lexer.New(tt.input)
			p := parser.New(l)

			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("wrong length. expected=%d, got=%d", 1, len(program.Statements))
			}

			stmt := program.Statements[0]
			if !testLetStatement(t, stmt, tt.expectedIdentifier) {
				return
			}

			val := stmt.(*ast.LetStatement).Value
			if !testLiteralExpression(t, val, tt.expectedValue) {
				return
			}
		})
	}
}

func TestReturnStatements(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			l := lexer.New(tt.input)
			p := parser.New(l)

			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("wrong length. expected=%d, got=%d", 1, len(program.Statements))
			}

			stmt := program.Statements[0]
			if !testReturnStatement(t, stmt) {
				return
			}

			ret := stmt.(*ast.ReturnStatement)
			if !testLiteralExpression(t, ret.ReturnValue, tt.expectedValue) {
				return
			}
		})
	}
}

func TestIdentifierExpression(t *testing.T) {
	t.Parallel()

	input := "foobar;"

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	testLiteralExpression(t, stmt.Expression, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	t.Parallel()

	input := "5;"

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	testLiteralExpression(t, stmt.Expression, 5)
}

func TestBooleanExpression(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			l := lexer.New(tt.input)
			p := parser.New(l)

			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program.Statements does not contain enough statements. expected=%d, got=%d", 1, len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statments[0] not *ast.ExpressionStatement. got=%T", stmt)
			}

			if !testLiteralExpression(t, stmt.Expression, tt.expectedBoolean) {
				return
			}
		})
	}
}

func TestPrefixExpression(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input         string
		operator      string
		expectedValue interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statement. got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not *ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.expectedValue) {
			return
		}
	}
}

func TestInfixExpression(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			l := lexer.New(tt.input)
			p := parser.New(l)

			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
			}

			infixExp, ok := stmt.Expression.(*ast.InfixExpression)
			if !ok {
				t.Fatalf("stmt.Expression not *ast.InfixExpression. got=%T", stmt.Expression)
			}

			if !testLiteralExpression(t, infixExp.Left, tt.leftValue) {
				return
			}

			if infixExp.Operator != tt.operator {
				t.Fatalf("infixExp.Operator not %s. got=%s", infixExp.Operator, tt.operator)
				return
			}

			if !testLiteralExpression(t, infixExp.Right, tt.rightValue) {
				return
			}
		})
	}
}

func TestIfExpression(t *testing.T) {
	t.Parallel()

	input := "if (x < y) { x }"

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("exp.Consequence does not contain %d statements. got=%d", 1, len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Fatalf("exp.Alternative is not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	t.Parallel()

	input := "if (x < y) { x } else { y }"

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("exp.Consequence does not contain %d statements. got=%d", 1, len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Alternative.Statements[0] is not *ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteral(t *testing.T) {
	t.Parallel()

	input := "fn(x, y) { x + y; }"

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("wrong length. expected=%d, got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("wrong type. expected=*ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("wrong type. expected=*ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("wrong length. expected=%d, got=%d", 2, len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("wrong length. expected=%d, got=%d", 1, len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("wrong type. expected=*ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			l := lexer.New(tt.input)
			p := parser.New(l)

			program := p.ParseProgram()
			checkParserErrors(t, p)

			stmt := program.Statements[0].(*ast.ExpressionStatement)
			function := stmt.Expression.(*ast.FunctionLiteral)

			if len(function.Parameters) != len(tt.expectedParams) {
				t.Errorf("parameters length wrong. expected=%d, got=%d", len(tt.expectedParams), len(function.Parameters))
			}

			for i, ident := range tt.expectedParams {
				testLiteralExpression(t, function.Parameters[i], ident)
			}
		})
	}
}

func TestCallExpression(t *testing.T) {
	t.Parallel()

	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("wrong statements length. expected=%d, got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("wrong statement type. expected=*ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	callExp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("wrong expression type. expected=*ast.CallExpression, got=%T", stmt.Expression)
	}

	if !testIdentifier(t, callExp.Function, "add") {
		return
	}

	if len(callExp.Arguments) != 3 {
		t.Fatalf("wrong arguments length. expected=%d, got=%d", 3, len(callExp.Arguments))
	}

	if !testLiteralExpression(t, callExp.Arguments[0], 1) {
		return
	}

	if !testInfixExpression(t, callExp.Arguments[1], 2, "*", 3) {
		return
	}

	if !testInfixExpression(t, callExp.Arguments[2], 4, "+", 5) {
		return
	}
}

func TestCallExpressionParameters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{"add();", "add", []string{}},
		{"add(1);", "add", []string{"1"}},
		{"add(1, 2 * 3, 4 + 5);", "add", []string{"1", "(2 * 3)", "(4 + 5)"}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			l := lexer.New(tt.input)
			p := parser.New(l)

			program := p.ParseProgram()
			checkParserErrors(t, p)

			stmt := program.Statements[0].(*ast.ExpressionStatement)
			callExp, ok := stmt.Expression.(*ast.CallExpression)
			if !ok {
				t.Fatalf("wrong expression type. expected=*ast.CallExpression, got=%T", stmt.Expression)
			}

			if !testIdentifier(t, callExp.Function, tt.expectedIdent) {
				return
			}

			if len(callExp.Arguments) != len(tt.expectedArgs) {
				t.Fatalf("wrong arguments length. expected=%d, got=%d", len(tt.expectedArgs), len(callExp.Arguments))
			}

			for i, arg := range tt.expectedArgs {
				if callExp.Arguments[i].String() != arg {
					t.Fatalf("wrong argument. expected=%s, got=%s", arg, callExp.Arguments[i].String())
				}
			}
		})
	}
}

func TestOperatorPrecedence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			l := lexer.New(tt.input)
			p := parser.New(l)

			program := p.ParseProgram()
			checkParserErrors(t, p)

			got := program.String()
			if got != tt.expected {
				t.Fatalf("expected=%q, got=%q", tt.expected, got)
			}
		})
	}
}
