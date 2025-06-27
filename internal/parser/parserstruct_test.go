package parser_test

import (
	"testing"
	"github.com/ithinkiborkedit/niftelv2.git/internal/lexer"
	"github.com/ithinkiborkedit/niftelv2.git/internal/parser"
	"github.com/ithinkiborkedit/niftelv2.git/internal/interpreter"
	"github.com/ithinkiborkedit/niftelv2.git/internal/value"
)

func TestStructCreationAndFields(t *testing.T) {
	src := `
struct Person {
    name: string
    age: int
}
var p: Person = Person { name: "Alice", age: 42 }
`
	lex := lexer.New(src)
	par := parser.New(lex)
	stmts, err := par.Parse()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	interp := interpreter.NewInterpreter()
	for _, stmt := range stmts {
		if err := interp.Execute(stmt); err != nil {
			t.Fatalf("Interpreter error: %v", err)
		}
	}
	// Check that variable exists and fields are set
	pVal, ok := interp.Globals["p"]
	if !ok {
		t.Fatalf("Variable 'p' not found in interpreter globals")
	}
	if pVal.Type != value.ValueStruct {
		t.Fatalf("Expected p to be a struct, got: %v", pVal.Type)
	}
	fields := pVal.Data.(map[string]value.Value)
	if fields["name"].String() != `"Alice"` {
		t.Fatalf("Expected p.name == \"Alice\", got %v", fields["name"].String())
	}
	if fields["age"].Data != float64(42) {
		t.Fatalf("Expected p.age == 42, got %v", fields["age"].Data)
	}
}

func TestStructWithMethodCall(t *testing.T) {
	src := `
struct Point {
    x: int
    y: int

    func sum() -> int {
        return x + y
    }
}
var p: Point = Point { x: 3, y: 4 }
var s: int = p.sum()
`
	lex := lexer.New(src)
	par := parser.New(lex)
	stmts, err := par.Parse()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	interp := interpreter.NewInterpreter()
	for _, stmt := range stmts {
		if err := interp.Execute(stmt); err != nil {
			t.Fatalf("Interpreter error: %v", err)
		}
	}
	sVal, ok := interp.Globals["s"]
	if !ok {
		t.Fatalf("Variable 's' not found in interpreter globals")
	}
	if sVal.Type != value.ValueNumber {
		t.Fatalf("Expected s to be a number, got %v", sVal.Type)
	}
	if sVal.Data != float64(7) {
		t.Fatalf("Expected s == 7, got %v", sVal.Data)
	}
}

func TestStructMethodOnLiteralInstance(t *testing.T) {
	src := `
struct Box {
    width: int
    height: int

    func area() -> int {
        return width * height
    }
}
var a: int = Box { width: 2, height: 3 }.area()
`
	lex := lexer.New(src)
	par := parser.New(lex)
	stmts, err := par.Parse()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	interp := interpreter.NewInterpreter()
	for _, stmt := range stmts {
		if err := interp.Execute(stmt); err != nil {
			t.Fatalf("Interpreter error: %v", err)
		}
	}
	aVal, ok := interp.Globals["a"]
	if !ok {
		t.Fatalf("Variable 'a' not found in interpreter globals")
	}
	if aVal.Type != value.ValueNumber {
		t.Fatalf("Expected a to be a number, got %v", aVal.Type)
	}
	if aVal.Data != float64(6) {
		t.Fatalf("Expected a == 6, got %v", aVal.Data)
	}
}