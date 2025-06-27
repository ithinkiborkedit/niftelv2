package parser_test

import (
	"testing"

	"github.com/ithinkiborkedit/niftelv2.git/internal/lexer"
	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	"github.com/ithinkiborkedit/niftelv2.git/internal/parser"
)

// Helper: returns a TokenSource from a string source.
func tokenSourceFromString(source string) lexer.TokenSource {
	return lexer.New(source)
}

func TestParser_Parse(t *testing.T) {
	source := `
		var x: int = 42
		var y: int = x + 1

		func foo() -> (int,int) {
		  return 1
		}
		print(x)
	`
	// Use the lexer directly as TokenSource
	tokSrc := tokenSourceFromString(source)
	p := parser.New(tokSrc)

	stmts, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(stmts) != 3 {
		t.Fatalf("Expected 3 statements, got %d", len(stmts))
	}

	// Check first statement: var x: int = 42
	var1, ok := stmts[0].(*ast.VarStmt)
	if !ok {
		t.Fatalf("First stmt is not VarStmt, got %T", stmts[0])
	}
	if var1.Names[0].Lexeme != "x" {
		t.Errorf("Expected var name 'x', got %q", var1.Names[0].Lexeme)
	}

	// Check second statement: var y: int = x + 1
	var2, ok := stmts[1].(*ast.VarStmt)
	if !ok {
		t.Fatalf("Second stmt is not VarStmt, got %T", stmts[1])
	}
	if var2.Names[0].Lexeme != "y" {
		t.Errorf("Expected var name 'y', got %q", var2.Names[0].Lexeme)
	}

	// Check third statement: print(x)
	printStmt, ok := stmts[2].(*ast.PrintStmt)
	if !ok {
		t.Fatalf("Third stmt is not PrintStmt, got %T", stmts[2])
	}
	_ = printStmt // Can add more assertions if desired
}
func TestParseMultipleReturnValues(t *testing.T) {
	src := `
func foo() -> (int, int) {
    return 1, 2
}
`
	lex := lexer.New(src)
	par := parser.New(lex)
	stmts, err := par.Parse()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}
	if len(stmts) == 0 {
		t.Fatalf("No statements parsed")
	}
	fn, ok := stmts[0].(*ast.FuncStmt)
	if !ok {
		t.Fatalf("First statement is not a FuncStmt")
	}
	found := false
	for _, stmt := range fn.Body.Statements {
		ret, ok := stmt.(*ast.ReturnStmt)
		if ok {
			found = true
			// --- Updated assertion: ---
			if len(ret.Values) != 2 {
				t.Fatalf("Return should have 2 values, got %d", len(ret.Values))
			}
			lit1, ok1 := ret.Values[0].(*ast.LiteralExpr)
			lit2, ok2 := ret.Values[1].(*ast.LiteralExpr)
			if !ok1 || !ok2 {
				t.Fatalf("Return values not parsed as literals")
			}
			if lit1.Value.Lexeme != "1" || lit2.Value.Lexeme != "2" {
				t.Fatalf("Return values not correct: got %q, %q", lit1.Value.Lexeme, lit2.Value.Lexeme)
			}
		}
	}
	if !found {
		t.Fatalf("Did not find a ReturnStmt in function body")
	}
}
