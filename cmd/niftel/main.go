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
	"github.com/ithinkiborkedit/niftelv2.git/internal/runtimecontrol"
	"github.com/ithinkiborkedit/niftelv2.git/internal/value"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("TOP-LEVEL PANIC: %#v\n", r)
			panic(r)
		}
	}()
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
		fmt.Printf("DEBUG stms: %#v\n", stmts)
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
				var result value.Value
				var evalErr error
				func() {
					defer func() {
						if r := recover(); r != nil {
							if ret, ok := r.(runtimecontrol.ReturnValue); ok {
								result = ret.Value
								evalErr = nil
							} else {
								panic(r)
							}
						}
					}()
					result, evalErr = interp.Eval(s.Expr)
				}()
				if evalErr != nil {
					fmt.Printf("Runtime error: %v\n", evalErr)
					break
				}
				if !result.IsNull() {
					fmt.Println(result.String())
				}
			}
		}
		buffer.Reset()
		prompt = ">>> "
	}
}
