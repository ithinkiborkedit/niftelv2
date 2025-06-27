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

// Helper: counts '{' and '}' in a line, ignoring those in strings
func countBraces(line string) (open, close int) {
	inString := false
	stringChar := byte(0)
	for i := 0; i < len(line); i++ {
		c := line[i]
		if inString {
			if c == stringChar {
				inString = false
			} else if c == '\\' && i+1 < len(line) {
				i++ // skip escaped char
			}
			continue
		}
		if c == '"' || c == '\'' {
			inString = true
			stringChar = c
			continue
		}
		if c == '{' {
			open++
		}
		if c == '}' {
			close++
		}
	}
	return
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("TOP-LEVEL PANIC: %#v\n", r)
			panic(r)
		}
	}()

	interp := interpreter.NewInterpreter()
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Niftel REPL v0")
	prompt := ">>> "

	for {
		var buffer strings.Builder
		openBraces := 0
		firstLine := true

		// Read one or more lines depending on { ... }
		for {
			fmt.Print(prompt)
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			if firstLine && strings.TrimSpace(line) == "" {
				// ignore blank lines at start
				continue
			}
			buffer.WriteString(line)

			o, c := countBraces(line)
			openBraces += o - c

			if openBraces > 0 || (firstLine && strings.Contains(line, "{")) {
				// inside block: switch to ... prompt
				prompt = "... "
			} else {
				break // block closed or no block at all
			}
			firstLine = false
		}

		// Debug: show what was read
		fmt.Printf("[REPL RAW BUFFER]\n%q\n", buffer.String())
		for i, c := range buffer.String() {
			fmt.Printf("%03d: %q (%d)\n", i, c, c)
		}

		lex := lexer.New(buffer.String())
		par := parser.New(lex)
		stmts, err := par.Parse()
		fmt.Printf("DEBUG stms: %#v\n", stmts)
		if err == parser.ErrIncomplete {
			prompt = "... "
			continue
		} else if err != nil {
			fmt.Printf("Parser error: %v \n", err)
			prompt = ">>> "
			continue
		}

		for _, stmt := range stmts {
			fmt.Printf("REPL: STATEMENT type: %T\n", stmt)
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
				fmt.Printf("DEBUG result: %#v\n", result)
				if !result.IsNull() {
					fmt.Println(result.String())
				}
			default:
				err := interp.Execute(stmt)
				if err != nil {
					fmt.Printf("Runtime error: %v\n", err)
				}
			}
		}
		// Reset prompt for next input
		prompt = ">>> "
	}
}