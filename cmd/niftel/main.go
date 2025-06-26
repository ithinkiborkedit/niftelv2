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
		// if line == "" {
		// 	continue
		// }

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

		// if err != nil {
		// 	fmt.Printf("Parse error %v\n", err)
		// 	continue
		// }

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
				// var result value.Value
				// var err error
				// func() {
				// 	defer func() {
				// 		if r := recover(); r != nil {
				// 			if ret, ok := r.(runtimecontrol.ReturnValue); ok {
				// 				result = ret.Value
				// 				err = nil
				// 			} else {
				// 				panic(r)
				// 			}
				// 		}
				// 	}()
				// 	result, err = interp.Eval(s.Expr)
				// }()
				// if err != nil {
				// 	fmt.Printf("Runtime error: %v\n", err)
				// 	break
				// }
				// fmt.Printf("DEBUG result: %#v\n", result)
				// if !result.IsNull() {
				// 	fmt.Println(result.String())
				// }
				// fmt.Println(result.String())
				// val, err := interp.Eval(s.Expr)
				// if err != nil {
				// 	fmt.Printf("Runtime error: %v\n", err)
				// 	break
				// }
				// fmt.Println(val.String())
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
