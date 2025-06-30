package environment

import (
	"fmt"

	"github.com/ithinkiborkedit/niftelv2.git/internal/symtable"
	"github.com/ithinkiborkedit/niftelv2.git/internal/value"
)

type Environment struct {
	symbols *symtable.SymbolTable
	values  map[string]value.Value
	// variables map[string]value.Value
	enclosing *Environment
	// types     map[string]*value.StructType
}

func NewEnvironment(parent *Environment) *Environment {
	var parentTable *symtable.SymbolTable
	if parent != nil {
		parentTable = parent.symbols

	}

	return &Environment{
		symbols:   symtable.NewSymbolTable(parentTable),
		values:    make(map[string]value.Value),
		enclosing: parent,
	}
}

func (e *Environment) envForSymbol(kind symtable.SymbolKind, name string) *Environment {
	if e.symbols.HasLocal(kind, name) {
		return e
	}
	if e.enclosing != nil {
		return e.enclosing.envForSymbol(kind, name)

	}
	return nil
}

func (e *Environment) Parent() *Environment {
	return e.enclosing
}

func (e *Environment) SymbolTable() *symtable.SymbolTable {
	return e.symbols
}

func (e *Environment) AssignVar(name string, val value.Value) error {
	env := e.envForSymbol(symtable.SymbolVar, name)
	if env == nil {
		return fmt.Errorf("undefined  variable '%s'", name)
	}
	env.values[name] = val
	return nil
}

func (e *Environment) DefineVar(sym *symtable.VarSymbol) error {
	return e.symbols.DefineValue(sym)
}

func (e *Environment) DefineFunc(sym *symtable.FuncSymbol) error {
	return e.symbols.DefineValue(sym)
}

func (e *Environment) DefineType(sym *symtable.TypeSymbol) error {
	return e.symbols.DefineValue(sym)
}

func (e *Environment) DefineTypeParam(sym *symtable.TypeParamSymbol) error {
	return e.symbols.DefineValue(sym)
}

func (e *Environment) GetVar(name string) (value.Value, error) {
	env := e.envForSymbol(symtable.SymbolVar, name)
	if env == nil {
		return value.Null(), fmt.Errorf("undefined variable '%s'", name)
	}
	val, ok := env.values[name]
	if !ok {
		return value.Null(), fmt.Errorf("uninitialised variable '%s'", name)
	}
	return val, nil
}

func (e *Environment) LookupVar(name string) (*symtable.VarSymbol, bool) {
	sym, ok := e.symbols.Lookup(symtable.SymbolVar, name)
	if !ok {
		return nil, false
	}
	varSym, ok := sym.(*symtable.VarSymbol)
	return varSym, ok
}

func (e *Environment) LookupFunc(name string) (*symtable.FuncSymbol, bool) {
	sym, ok := e.symbols.Lookup(symtable.SymbolFuncs, name)
	if !ok {
		return nil, false
	}
	funcSym, ok := sym.(*symtable.FuncSymbol)
	return funcSym, ok
}

func (e *Environment) LookupType(name string) (*symtable.TypeSymbol, bool) {
	sym, ok := e.symbols.Lookup(symtable.SymbolTypes, name)
	if !ok {
		return nil, false
	}
	typeSym, ok := sym.(*symtable.TypeSymbol)
	return typeSym, ok
}

func (e *Environment) LookupTypeParam(name string) (*symtable.TypeParamSymbol, bool) {
	sym, ok := e.symbols.Lookup(symtable.SymbolTypeParams, name)
	if !ok {
		return nil, false
	}
	typeparamSym, ok := sym.(*symtable.TypeParamSymbol)
	return typeparamSym, ok
}

func (e *Environment) HasLocalVar(name string) bool {
	return e.symbols.HasLocal(symtable.SymbolVar, name)
}

func (e *Environment) HasLocalFunc(name string) bool {
	return e.symbols.HasLocal(symtable.SymbolFuncs, name)
}

func (e *Environment) HasLocalType(name string) bool {
	return e.symbols.HasLocal(symtable.SymbolTypes, name)
}

func (e *Environment) HasLocalTypeparam(name string) bool {
	return e.symbols.HasLocal(symtable.SymbolTypeParams, name)
}
