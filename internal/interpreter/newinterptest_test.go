package interpreter_test

import (
	"testing"

	"github.com/ithinkiborkedit/niftelv2.git/internal/environment"
	"github.com/ithinkiborkedit/niftelv2.git/internal/interpreter"
	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	token "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"
	"github.com/ithinkiborkedit/niftelv2.git/internal/symtable"
)

func TestInterpreter_VarDeclareAssignAndFuncCall(t *testing.T) {
	interp := interpreter.NewInterpreter()

	// 1. var x = 10
	stmtVar := &ast.VarStmt{
		Names: []token.Token{{Lexeme: "x"}},
		Init: &ast.LiteralExpr{
			Value: token.Token{Type: token.TokenNumber, Data: 10},
		},
	}

	res := interp.Execute(stmtVar)
	if res.Err != nil {
		t.Fatalf("var declaration failed: %v", res.Err)
	}
	got, err := interp.GetEnv().GetVar("x")
	if err != nil || got.Data.(float64) != 10 {
		t.Errorf("expected x == 10, got %v, err=%v", got, err)
	}

	// 2. x = 15
	stmtAssign := &ast.AssignStmt{
		Name:  token.Token{Lexeme: "x"},
		Value: &ast.LiteralExpr{Value: token.Token{Type: token.TokenNumber, Data: 15}},
	}
	res = interp.Execute(stmtAssign)
	if res.Err != nil {
		t.Fatalf("assign failed: %v", res.Err)
	}
	got, err = interp.GetEnv().GetVar("x")
	if err != nil || got.Data.(float64) != 15 {
		t.Errorf("expected x == 15, got %v, err=%v", got, err)
	}

	// 3. func foo() { x = 42 }
	stmtFunc := &ast.FuncStmt{
		Name:   token.Token{Lexeme: "foo"},
		Params: nil,
		Body: &ast.BlockStmt{
			Statements: []ast.Stmt{
				&ast.AssignStmt{
					Name:  token.Token{Lexeme: "x"},
					Value: &ast.LiteralExpr{Value: token.Token{Type: token.TokenNumber, Data: 42}},
				},
			},
		},
		Func: token.Token{Lexeme: "func"}, // Or whatever is correct for your AST
	}
	res = interp.Execute(stmtFunc)
	if res.Err != nil {
		t.Fatalf("func definition failed: %v", res.Err)
	}

	// 4. Call foo()
	callExpr := &ast.CallExpr{
		Callee: &ast.VariableExpr{Name: token.Token{Lexeme: "foo"}},
	}
	res = interp.Eval(callExpr)
	if res.Err != nil {
		t.Fatalf("function call failed: %v", res.Err)
	}
	got, err = interp.GetEnv().GetVar("x")
	if err != nil || got.Data.(float64) != 42 {
		t.Errorf("expected x == 42 after foo(), got %v, err=%v", got, err)
	}

	// 5. Shadowing: let’s try a new block
	block := &ast.BlockStmt{
		Statements: []ast.Stmt{
			&ast.VarStmt{
				Names: []token.Token{{Lexeme: "x"}},
				Init:  &ast.LiteralExpr{Value: token.Token{Type: token.TokenNumber, Data: 99}},
			},
		},
	}
	interp.PushEnv(environment.NewEnvironment(interp.GetEnv()))
	interp.ExecuteBlock(block, interp.GetEnv())
	got, err = interp.GetEnv().GetVar("x")
	if err != nil || got.Data.(float64) != 99 {
		t.Errorf("expected x == 99 in inner block, got %v, err=%v", got, err)
	}
	interp.PopEnv()

	got, err = interp.GetEnv().GetVar("x")
	if err != nil || got.Data.(float64) != 42 {
		t.Errorf("expected x == 42 after block, got %v, err=%v", got, err)
	}

	// 6. Error on redeclaring variable in same scope
	err = interp.GetEnv().DefineVar(&symtable.VarSymbol{
		SymName: "x",
		SymKind: symtable.SymbolVar,
		Mutable: true,
	})
	if err == nil {
		t.Errorf("expected error on redeclaring variable in same scope, got nil")
	}
}
