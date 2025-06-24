package interpreter_test

import (
	"testing"

	"github.com/ithinkiborkedit/niftelv2.git/internal/interpreter"
	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	token "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"
	"github.com/ithinkiborkedit/niftelv2.git/internal/value"
)

func TestInterpreter_EvaluateLiteralInt(t *testing.T) {
	interp := interpreter.NewInterpreter()

	intToken := token.Token{
		Type:   token.TokenNumber,
		Lexeme: "42",
		Data:   float64(42),
		Line:   1,
		Column: 1,
	}

	literalExpr := &ast.LiteralExpr{
		Value: intToken,
	}

	val, err := interp.Evaluate(literalExpr)
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}

	if val.Type != value.ValueInt {
		t.Fatalf("expected valueInt type got %v", val.Type)
	}

	if intVal, ok := val.Data.(float64); !ok || intVal != 42 {
		t.Errorf("expected int 42, got %v (ok=%v)", val.Data, ok)
	}
}

func TestInterpreter_EvaluateBinaryExprIntComparisson(t *testing.T) {
	interp := interpreter.NewInterpreter()

	intTokenL := token.Token{
		Type:   token.TokenNumber,
		Lexeme: "42",
		Data:   42,
		Line:   1,
		Column: 1,
	}

	intTokenR := token.Token{
		Type:   token.TokenNumber,
		Lexeme: "100",
		Data:   100,
		Line:   1,
		Column: 1,
	}

	left := &ast.LiteralExpr{Value: intTokenL}
	right := &ast.LiteralExpr{Value: intTokenR}

	binExpr := &ast.BinaryExpr{
		Left:     left,
		Operator: token.Token{Type: token.TokenLess, Lexeme: "<", Line: 1, Column: 2},
		Right:    right,
	}

	result, err := interp.Evaluate(binExpr)
	if err != nil {
		t.Fatalf("error evaluating binary expr: %v", err)
	}

	if result.Type != value.ValueBool {
		t.Fatalf("expected bool result, got %v", result.Type)
	}

	if val, ok := result.Data.(bool); !ok || val != true {
		t.Fatalf("expected true for 42 < 100, got %v", result.Data)
	}

	// literalExpr := &ast.LiteralExpr{
	// 	Value: intToken,
	// }

	// val, err := interp.Evaluate(literalExpr)
	// if err != nil {
	// 	t.Fatalf("Evaluate returned error: %v", err)
	// }

	// if val.Type != value.ValueInt {
	// 	t.Fatalf("expected valueInt type got %v", val.Type)
	// }

	// if intVal, ok := val.Data.(int); !ok || intVal != 42 {
	// 	t.Errorf("expected int 42, got %v (ok=%v)", val.Data, ok)
	// }
}
