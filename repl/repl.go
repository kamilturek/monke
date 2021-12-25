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
		fmt.Print(PROMPT)
		if !scanner.Scan() {
			return
		}

		l := lexer.New(scanner.Text())

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%v\n", tok)
		}
	}
}
