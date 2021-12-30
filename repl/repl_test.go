package repl_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kamilturek/monke/repl"
)

func TestStart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"x", ">> {Type:IDENT Literal:x}\n>> "},
		{
			"let x = 5;",
			strings.Join(
				[]string{
					">> {Type:LET Literal:let}",
					"{Type:IDENT Literal:x}",
					"{Type:= Literal:=}",
					"{Type:INT Literal:5}",
					"{Type:; Literal:;}",
					">> ",
				},
				"\n",
			),
		},
		{
			"5 + 10 / 5 * 2",
			strings.Join(
				[]string{
					">> {Type:INT Literal:5}",
					"{Type:+ Literal:+}",
					"{Type:INT Literal:10}",
					"{Type:/ Literal:/}",
					"{Type:INT Literal:5}",
					"{Type:* Literal:*}",
					"{Type:INT Literal:2}",
					">> ",
				},
				"\n",
			),
		},
	}

	for i, tt := range tests {
		in := bytes.NewBufferString(tt.input)
		out := &bytes.Buffer{}
		repl.Start(in, out)

		if out.String() != tt.expectedOutput {
			t.Fatalf("tests[%d] - output wrong. expected=%q, got=%q", i, tt.expectedOutput, out.String())
		}
	}
}
