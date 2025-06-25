package parser_test

import (
	"testing"

	"github.com/ithinkiborkedit/niftelv2.git/internal/lexer"
	"github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
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
	var1, ok := stmts[0].(*nifast.VarStmt)
	if !ok {
		t.Fatalf("First stmt is not VarStmt, got %T", stmts[0])
	}
	if var1.Name.Lexeme != "x" {
		t.Errorf("Expected var name 'x', got %q", var1.Name.Lexeme)
	}

	// Check second statement: var y: int = x + 1
	var2, ok := stmts[1].(*nifast.VarStmt)
	if !ok {
		t.Fatalf("Second stmt is not VarStmt, got %T", stmts[1])
	}
	if var2.Name.Lexeme != "y" {
		t.Errorf("Expected var name 'y', got %q", var2.Name.Lexeme)
	}

	// Check third statement: print(x)
	printStmt, ok := stmts[2].(*nifast.PrintStmt)
	if !ok {
		t.Fatalf("Third stmt is not PrintStmt, got %T", stmts[2])
	}
	_ = printStmt // Can add more assertions if desired
}
