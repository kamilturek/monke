package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/kamilturek/monke/lexer"
	"github.com/kamilturek/monke/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		_, err := io.WriteString(out, PROMPT)
		if err != nil {
			panic(err)
		}

		if !scanner.Scan() {
			return
		}

		l := lexer.New(scanner.Text())

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			tokFmt := fmt.Sprintf("%+v\n", tok)
			if _, err := io.WriteString(out, tokFmt); err != nil {
				panic(err)
			}
		}
	}
}
