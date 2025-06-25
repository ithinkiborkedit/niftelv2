package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ithinkiborkedit/niftelv2.git/internal/interpreter"
	"github.com/ithinkiborkedit/niftelv2.git/internal/lexer"
	"github.com/ithinkiborkedit/niftelv2.git/internal/parser"
)

func main() {
	interp := interpreter.NewInterpreter()
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Niftel REPL v0")
	for {
		fmt.Print(">>> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		lex := lexer.New(line)
		par := parser.New(lex)

		stmts, err := par.Parse()
		if err != nil {
			fmt.Printf("Parse error %v\n", err)
			continue
		}

		err = interp.Execute(stmts)

		// node, err := parser
	}
}
