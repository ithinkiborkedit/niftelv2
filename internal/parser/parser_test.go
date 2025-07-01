package parser_test

import (
	"testing"

	"github.com/ithinkiborkedit/niftelv2.git/internal/lexer"
	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	"github.com/ithinkiborkedit/niftelv2.git/internal/parser"
)

// Test var x: int = 10
func TestParseVarDeclarationWithType(t *testing.T) {
	src := "var x: int = 10"
	lex := lexer.New(src)
	par := parser.New(lex)

	stmts, err := par.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
	varStmt, ok := stmts[0].(*ast.VarStmt)
	if !ok {
		t.Fatalf("expected VarStmt, got %T", stmts[0])
	}
	if len(varStmt.Names) != 1 || varStmt.Names[0].Lexeme != "x" {
		t.Errorf("expected var name 'x', got %+v", varStmt.Names)
	}
	if varStmt.Type.Lexeme != "int" {
		t.Errorf("expected type 'int', got '%s'", varStmt.Type.Lexeme)
	}
}

// Test func add(a: int, b: int) -> int { return a + b }
func TestParseFuncDeclaration(t *testing.T) {
	src := `func add(a: int, b: int) -> int {
		return a + b
	}`
	lex := lexer.New(src)
	par := parser.New(lex)

	stmts, err := par.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
	funcStmt, ok := stmts[0].(*ast.FuncStmt)
	if !ok {
		t.Fatalf("expected FuncStmt, got %T", stmts[0])
	}
	if funcStmt.Name.Lexeme != "add" {
		t.Errorf("expected func name 'add', got '%s'", funcStmt.Name.Lexeme)
	}
	if len(funcStmt.Params) != 2 {
		t.Errorf("expected 2 params, got %d", len(funcStmt.Params))
	}
	if funcStmt.Params[0].Name.Lexeme != "a" || funcStmt.Params[0].Type.Lexeme != "int" {
		t.Errorf("expected param 'a: int', got '%s: %s'", funcStmt.Params[0].Name.Lexeme, funcStmt.Params[0].Type.Lexeme)
	}
	if funcStmt.Params[1].Name.Lexeme != "b" || funcStmt.Params[1].Type.Lexeme != "int" {
		t.Errorf("expected param 'b: int', got '%s: %s'", funcStmt.Params[1].Name.Lexeme, funcStmt.Params[1].Type.Lexeme)
	}
	if len(funcStmt.Return) != 1 || funcStmt.Return[0].Lexeme != "int" {
		t.Errorf("expected return type 'int', got %+v", funcStmt.Return)
	}
}

// Test struct Point { x: int y: int }
func TestParseStructDeclaration(t *testing.T) {
	src := `
struct Point {
	x: int
	y: int
}
`
	lex := lexer.New(src)
	par := parser.New(lex)

	stmts, err := par.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
	structStmt, ok := stmts[0].(*ast.StructStmt)
	if !ok {
		t.Fatalf("expected StructStmt, got %T", stmts[0])
	}
	if structStmt.Name.Lexeme != "Point" {
		t.Errorf("expected struct name 'Point', got '%s'", structStmt.Name.Lexeme)
	}
	if len(structStmt.Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(structStmt.Fields))
	}
	if structStmt.Fields[0].Names[0].Lexeme != "x" || structStmt.Fields[0].Type.Lexeme != "int" {
		t.Errorf("expected field 'x: int', got '%s: %s'", structStmt.Fields[0].Names[0].Lexeme, structStmt.Fields[0].Type.Lexeme)
	}
	if structStmt.Fields[1].Names[0].Lexeme != "y" || structStmt.Fields[1].Type.Lexeme != "int" {
		t.Errorf("expected field 'y: int', got '%s: %s'", structStmt.Fields[1].Names[0].Lexeme, structStmt.Fields[1].Type.Lexeme)
	}
}
