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
		{"x", ">> {IDENT x}\n>> "},
		{"let x = 5;", ">> {LET let}\n{IDENT x}\n{= =}\n{INT 5}\n{; ;}\n>> "},
		{"5 + 10 / 5 * 2", ">> {INT 5}\n{+ +}\n{INT 10}\n{/ /}\n{INT 5}\n{* *}\n{INT 2}\n>> "},
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
