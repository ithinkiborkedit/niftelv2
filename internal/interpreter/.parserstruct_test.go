package interpreter_test

import (
	"testing"

	"github.com/ithinkiborkedit/niftelv2.git/internal/interpreter"
	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	token "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"
	"github.com/ithinkiborkedit/niftelv2.git/internal/value"
)

func TestStructDeclarationAndInstantiation(t *testing.T) {
	value.ResetTypes()
	value.BuiltinTypesInit()
	interp := interpreter.NewInterpreter()

	// 1. Declare struct Point { x: int, y: int }
	pointStruct := &ast.StructStmt{
		Name: token.Token{Lexeme: "Point"},
		Fields: []ast.VarStmt{
			{
				Name: token.Token{Lexeme: "x"},
				Type: token.Token{Lexeme: "int"},
			},
			{
				Name: token.Token{Lexeme: "y"},
				Type: token.Token{Lexeme: "int"},
			},
		},
	}

	// 2. Register the struct type
	if err := interp.Execute(pointStruct); err != nil {
		t.Fatalf("struct registration failed: %v", err)
	}

	// 3. Instantiate struct: var p: Point = Point { x: 2, y: 3 }
	structLiteral := &ast.StructLiteralExpr{
		TypeName: token.Token{Lexeme: "Point"},
		Fields: map[string]ast.Expr{
			"x": &ast.LiteralExpr{Value: token.Token{Type: token.TokenNumber, Data: float64(2)}},
			"y": &ast.LiteralExpr{Value: token.Token{Type: token.TokenNumber, Data: float64(3)}},
		},
	}

	varStmt := &ast.VarStmt{
		Name: token.Token{Lexeme: "p"},
		Type: token.Token{Lexeme: "Point"},
		Init: structLiteral,
	}

	if err := interp.Execute(varStmt); err != nil {
		t.Fatalf("struct instantiation failed: %v", err)
	}

	// 4. Access fields on p
	val, err := interp.GetEnv().Get("p")
	if err != nil {
		t.Fatalf("failed to get variable 'p': %v", err)
	}

	// Check value type
	if val.Type != value.ValueStruct {
		t.Fatalf("expected ValueStruct, got %v", val.Type)
	}

	pointInstance, ok := val.Data.(*value.StructInstance)
	if !ok {
		t.Fatalf("expected *StructInstance, got %T", val.Data)
	}

	if got := pointInstance.Fields["x"]; got.Data.(float64) != 2 {
		t.Errorf("expected x=2, got %v", got.Data)
	}
	if got := pointInstance.Fields["y"]; got.Data.(float64) != 3 {
		t.Errorf("expected y=3, got %v", got.Data)
	}
}

func TestStructFieldDefaultsAndMissingFields(t *testing.T) {
	value.ResetTypes()
	value.BuiltinTypesInit()
	interp := interpreter.NewInterpreter()

	pointStruct := &ast.StructStmt{
		Name: token.Token{Lexeme: "Point"},
		Fields: []ast.VarStmt{
			{
				Name: token.Token{Lexeme: "x"},
				Type: token.Token{Lexeme: "int"},
			},
			{
				Name: token.Token{Lexeme: "y"},
				Type: token.Token{Lexeme: "int"},
			},
		},
	}
	if err := interp.Execute(pointStruct); err != nil {
		t.Fatalf("struct registration failed: %v", err)
	}

	// Leave out "y" field in literal, should be filled with Null()
	structLiteral := &ast.StructLiteralExpr{
		TypeName: token.Token{Lexeme: "Point"},
		Fields: map[string]ast.Expr{
			"x": &ast.LiteralExpr{Value: token.Token{Type: token.TokenNumber, Data: float64(7)}},
			// "y" is missing
		},
	}

	varStmt := &ast.VarStmt{
		Name: token.Token{Lexeme: "p2"},
		Type: token.Token{Lexeme: "Point"},
		Init: structLiteral,
	}
	if err := interp.Execute(varStmt); err != nil {
		t.Fatalf("struct instantiation failed: %v", err)
	}

	val, err := interp.GetEnv().Get("p2")
	if err != nil {
		t.Fatalf("failed to get variable 'p2': %v", err)
	}
	pointInstance, ok := val.Data.(*value.StructInstance)
	if !ok {
		t.Fatalf("expected *StructInstance, got %T", val.Data)
	}

	if got := pointInstance.Fields["x"]; got.Data.(float64) != 7 {
		t.Errorf("expected x=7, got %v", got.Data)
	}
	if got := pointInstance.Fields["y"]; !got.IsNull() {
		t.Errorf("expected y=null, got %v", got)
	}
}
