package symtable_test

import (
	"testing"

	"github.com/ithinkiborkedit/niftelv2.git/internal/symtable"
)

func TestSymbolTable_Basic(t *testing.T) {
	// Create root table (global scope)
	root := symtable.NewSymbolTable(nil)

	// Add a variable
	varSym := &symtable.VarSymbol{SymName: "x", SymKind: symtable.SymbolVar, Mutable: true}
	if err := root.DefineValue(varSym); err != nil {
		t.Fatalf("Failed to define variable: %v", err)
	}

	// Add a function (same name allowed in separate namespace)
	funcSym := &symtable.FuncSymbol{SymName: "x"}
	if err := root.DefineValue(funcSym); err != nil {
		t.Fatalf("Failed to define function: %v", err)
	}

	// Add a type (same name, type namespace)
	typeSym := &symtable.TypeSymbol{SymName: "x", SymKind: symtable.SymbolTypes}
	if err := root.DefineValue(typeSym); err != nil {
		t.Fatalf("Failed to define type: %v", err)
	}

	// Confirm separate namespaces
	if _, ok := root.Lookup(symtable.SymbolVar, "x"); !ok {
		t.Errorf("Variable x not found")
	}
	if _, ok := root.Lookup(symtable.SymbolFuncs, "x"); !ok {
		t.Errorf("Function x not found")
	}
	if _, ok := root.Lookup(symtable.SymbolTypes, "x"); !ok {
		t.Errorf("Type x not found")
	}
}

func TestSymbolTable_DuplicateError(t *testing.T) {
	root := symtable.NewSymbolTable(nil)
	varSym := &symtable.VarSymbol{SymName: "foo", SymKind: symtable.SymbolVar}
	if err := root.DefineValue(varSym); err != nil {
		t.Fatalf("Failed to define var foo: %v", err)
	}
	dupSym := &symtable.VarSymbol{SymName: "foo", SymKind: symtable.SymbolVar}
	if err := root.DefineValue(dupSym); err == nil {
		t.Errorf("Expected error for duplicate variable, got nil")
	}
}

func TestSymbolTable_ParentLookup(t *testing.T) {
	root := symtable.NewSymbolTable(nil)
	child := symtable.NewSymbolTable(root)

	varSym := &symtable.VarSymbol{SymName: "a", SymKind: symtable.SymbolVar}
	if err := root.DefineValue(varSym); err != nil {
		t.Fatalf("Failed to define var in root: %v", err)
	}
	// Should find in parent
	if _, ok := child.Lookup(symtable.SymbolVar, "a"); !ok {
		t.Errorf("Child failed to find var from parent scope")
	}
	// Now shadow in child
	shadow := &symtable.VarSymbol{SymName: "a", SymKind: symtable.SymbolVar}
	if err := child.DefineValue(shadow); err != nil {
		t.Fatalf("Failed to define var in child: %v", err)
	}
	sym, _ := child.Lookup(symtable.SymbolVar, "a")
	if sym != shadow {
		t.Errorf("Lookup did not return shadowed symbol from child")
	}
}

func TestSymbolTable_TypeParams(t *testing.T) {
	root := symtable.NewSymbolTable(nil)
	tp := &symtable.TypeParamSymbol{Symnam: "T"}
	if err := root.DefineValue(tp); err != nil {
		t.Fatalf("Failed to define type param: %v", err)
	}
	if _, ok := root.Lookup(symtable.SymbolTypeParams, "T"); !ok {
		t.Errorf("Type param T not found in table")
	}
}