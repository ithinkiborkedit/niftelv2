package parser_test

import (
	"strings"
	"testing"

	"github.com/ithinkiborkedit/niftelv2.git/internal/lexer"
	"github.com/ithinkiborkedit/niftelv2.git/internal/parser"
	token "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"
	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
)

func TestParser_MultipleVarDeclaration(t *testing.T) {
	src := "var a, b = pair(1, 2)"
	lex := lexer.NewLexer(strings.NewReader(src))
	p := parser.New(lex)

	stmts, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(stmts))
	}

	vs, ok := stmts[0].(*ast.VarStmt)
	if !ok {
		t.Fatalf("Statement not VarStmt: %T", stmts[0])
	}

	if len(vs.Names) != 2 {
		t.Errorf("Expected 2 variable names, got %d", len(vs.Names))
	}
	if vs.Names[0].Lexeme != "a" || vs.Names[1].Lexeme != "b" {
		t.Errorf("Unexpected names: %s, %s", vs.Names[0].Lexeme, vs.Names[1].Lexeme)
	}
	if vs.Type.Lexeme != "" {
		t.Errorf("Expected no type annotation, got %s", vs.Type.Lexeme)
	}

	call, ok := vs.Init.(*ast.CallExpr)
	if !ok {
		t.Fatalf("Init is not a CallExpr: %T", vs.Init)
	}
	if varExpr, ok := call.Callee.(*ast.VariableExpr); !ok || varExpr.Name.Lexeme != "pair" {
		t.Errorf("Call callee is not 'pair': %v", call.Callee)
	}
}