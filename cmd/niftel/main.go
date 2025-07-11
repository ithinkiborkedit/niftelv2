package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ithinkiborkedit/niftelv2.git/internal/codegen"
	"github.com/ithinkiborkedit/niftelv2.git/internal/interpreter"
	"github.com/ithinkiborkedit/niftelv2.git/internal/lexer"
	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	"github.com/ithinkiborkedit/niftelv2.git/internal/parser"
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

func compileProject(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read file: %v\n", err)
		os.Exit(2)
	}
	source := string(data)

	lex := lexer.New(source)
	par := parser.New(lex)
	stmts, err := par.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		os.Exit(3)
	}

	cg := codegen.NewCodeGen()
	ir, err := cg.GenerateLLVM(stmts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "code gen error: %v\n", err)
		os.Exit(4)
	}

	baseName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	llFile := baseName + ".ll"
	objFile := baseName + ".o"
	nifFile := baseName + ""

	if err := os.WriteFile(llFile, []byte(ir), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write LLVM IR: %v\n", err)
		os.Exit(5)
	}

	llcCmd := exec.Command("llc", "-march=arm64", llFile, "-filetype=obj", "-o", objFile)
	llcCmd.Stdout = os.Stdout
	llcCmd.Stderr = os.Stderr
	if err := llcCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "llc failed: %v\n", err)
		os.Exit(6)
	}

	clangCmd := exec.Command("clang", objFile, "-o", nifFile)
	clangCmd.Stdout = os.Stdout
	clangCmd.Stderr = os.Stderr
	if err := clangCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "clang failed: %v\n", err)
		os.Exit(7)
	}
	// fmt.Printf("LLVM IR WRITTEN to: %s\n", outfile)
}

func runFile(path string, interp *interpreter.Interpreter) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file: %v\n", path)
		os.Exit(2)
	}
	reader := bufio.NewReader(strings.NewReader(string(data)))

	for {
		var buffer strings.Builder
		openBraces := 0
		firstLine := true

		for {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "Read error: %v\n", err)
				return
			}
			if err == io.EOF && line == "" && buffer.Len() == 0 {
				return
			}
			if firstLine && strings.TrimSpace(line) == "" {
				if err == io.EOF {
					return
				}
				continue
			}
			buffer.WriteString(line)

			o, c := countBraces(line)
			openBraces += o - c

			if openBraces > 0 || (firstLine && strings.Contains(line, "{")) {

			} else {
				break
			}
			firstLine = false
			if err == io.EOF {
				break
			}
		}
		if buffer.Len() == 0 {
			return
		}
		lex := lexer.New(buffer.String())
		par := parser.New(lex)
		stmts, err := par.Parse()
		if err == parser.ErrIncomplete {
			continue
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "parser error: %v\n", err)
			continue
		}
		for _, stmt := range stmts {
			result := interp.Execute(stmt)
			if result.Err != nil {
				fmt.Fprintf(os.Stderr, "runtime error: %v\n", result.Err)
			}
		}
	}
}

func main() {
	value.BuiltinTypesInit()
	interp := interpreter.NewInterpreter()
	// value.BuiltinTypesInit()
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "compile":
			if len(os.Args) < 3 {
				fmt.Fprintf(os.Stderr, "Usgae %s compile <source-code-file.nif>\n", os.Args[0])
				os.Exit(1)
			}
			compileProject(os.Args[2])
			return
		default:
			interp.ShouldPrintResults = false
			runFile(os.Args[1], interp)
			return
		}
	}
	interp.ShouldPrintResults = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("TOP-LEVEL PANIC: %#v\n", r)
			panic(r)
		}
	}()

	interp.ShouldPrintResults = true
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
		// fmt.Printf("[REPL RAW BUFFER]\n%q\n", buffer.String())
		for i, c := range buffer.String() {
			fmt.Printf("%03d: %q (%d)\n", i, c, c)
		}

		lex := lexer.New(buffer.String())
		par := parser.New(lex)
		stmts, err := par.Parse()

		// fmt.Printf("DEBUG stms: %#v\n", stmts)
		if err == parser.ErrIncomplete {
			prompt = "... "
			continue
		} else if err != nil {
			fmt.Printf("Parser error: %v \n", err)
			prompt = ">>> "
			continue
		}

		for _, stmt := range stmts {
			// fmt.Printf("REPL: STATEMENT type: %T\n", stmt)
			switch s := stmt.(type) {
			case *ast.ExprStmt:
				res := interp.Eval(s.Expr)
				if res.Err != nil {
					fmt.Printf("Runtime error: %v\n", res.Err)
					break
				}
				result := res.Value

				if !result.IsNull() {
					fmt.Println(result.String())
				}
			default:
				result := interp.Execute(stmt)
				if result.Err != nil {
					fmt.Printf("Runtime Error %v\n", result.Err)
				}
				// err := interp.Execute(stmt)
				// if err != nil {
				// 	fmt.Printf("Runtime error: %v\n", err)
				// }
			}
		}
		// Reset prompt for next input
		prompt = ">>> "
	}
}
