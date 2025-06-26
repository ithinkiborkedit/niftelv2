package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ithinkiborkedit/niftelv2.git/internal/interpreter"
	"github.com/ithinkiborkedit/niftelv2.git/internal/lexer"
	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	"github.com/ithinkiborkedit/niftelv2.git/internal/parser"
)

func main() {
	var buffer strings.Builder
	interp := interpreter.NewInterpreter()
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Niftel REPL v0")
	prompt := ">>> "
	for {
		fmt.Print(prompt)
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.TrimSpace(line) == "" {
			continue
		}

		buffer.WriteString(line)

		lex := lexer.New(buffer.String())
		par := parser.New(lex)
		stmts, err := par.Parse()

		if err == parser.ErrIncomplete {
			prompt = "..."
			continue
		} else if err != nil {
			fmt.Printf("Parser error: %v \n", err)
			buffer.Reset()
			continue
		}

		for _, stmt := range stmts {
			switch s := stmt.(type) {
			case *ast.ExprStmt:
				result, err := interp.Eval(s.Expr)
				if err != nil {
					fmt.Printf("Runtime error: %v\n", err)
					break
				}
				fmt.Printf("DEBGUG result: %#v\n", result)
				if !result.IsNull() {
					fmt.Println(result.String())
				}
			default:
				err := interp.Execute(stmt)
				if err != nil {
					fmt.Printf("Runtime error: %v\n", err)
					break
				}
			}
		}
		buffer.Reset()
		prompt = ">>> "
	}
}
