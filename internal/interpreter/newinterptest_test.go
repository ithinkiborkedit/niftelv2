package interpreter_test

import (
	"testing"
	"fmt"

	"github.com/ithinkiborkedit/niftelv2.git/internal/interpreter"
	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	token "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"
)

// Helper: create a dummy FuncStmt node
func makeFuncStmt(name string) *ast.FuncStmt {
	return &ast.FuncStmt{
		Name: token.Token{Lexeme: name},
		Params: nil,
		Body:   &ast.BlockStmt{Statements: nil},
		Func:   token.Token{Lexeme: "func", Line: 1, Column: 1},
	}
}

func TestFuncStmtDuplicateError(t *testing.T) {
	interp := interpreter.NewInterpreter()

	// First function declaration should succeed
	funcA := makeFuncStmt("foo")
	result1 := interp.VisitFuncStmt(funcA)
	if result1.Err != nil {
		t.Fatalf("expected no error, got: %v", result1.Err)
	}

	// Second declaration with the same name in the same env should fail
	funcB := makeFuncStmt("foo")
	result2 := interp.VisitFuncStmt(funcB)
	if result2.Err == nil {
		t.Fatal("expected error for duplicate function, got nil")
	}
	want := "function 'foo' already defined in this scope"
	if result2.Err.Error() != want {
		t.Fatalf("expected error %q, got %q", want, result2.Err.Error())
	}
	fmt.Println("Duplicate function test passed:", result2.Err)
}