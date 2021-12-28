package repl

import (
	"bytes"
	"testing"
)

func TestStart(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"x", ">> {Type:IDENT Literal:x}\n>> "},
		{"let x = 5;", ">> {Type:LET Literal:let}\n{Type:IDENT Literal:x}\n{Type:= Literal:=}\n{Type:INT Literal:5}\n{Type:; Literal:;}\n>> "},
		{"5 + 10 / 5 * 2", ">> {Type:INT Literal:5}\n{Type:+ Literal:+}\n{Type:INT Literal:10}\n{Type:/ Literal:/}\n{Type:INT Literal:5}\n{Type:* Literal:*}\n{Type:INT Literal:2}\n>> "},
	}

	for i, tt := range tests {
		in := bytes.NewBufferString(tt.input)
		out := &bytes.Buffer{}
		Start(in, out)
		if out.String() != tt.expectedOutput {
			t.Fatalf("tests[%d] - output wrong. expected=%q, got=%q", i, tt.expectedOutput, out.String())
		}
	}
}
