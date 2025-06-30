package environment

import (
	"github.com/ithinkiborkedit/niftelv2.git/internal/symtable"
)

type Environment struct {
	symbols *symtable.SymbolTable
	// variables map[string]value.Value
	enclosing *Environment
	// types     map[string]*value.StructType
}

func NewEnvironment(parent *Environment) *Environment {
	var parentTable *symtable.SymbolTable
	if parent != nil {
		parentTable = parent.symbols

	}
	// return &Environment{
	// variables: make(map[string]value.Value),
	// 	enclosing: enclosing,
	// 	types:     make(map[string]*value.StructType),
	// }
	return &Environment{
		symbols:   symtable.NewSymbolTable(parentTable),
		enclosing: parent,
	}
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

//	func (e *Environment) Define(name string, val value.Value) {
//		fmt.Printf("ENV DEFINING '%s' as type %v", name, val.Type)
//		e.variables[name] = val
//	}
// func (e *Environment) Assign(name string, val value.Value) error {

// 	if _, ok := e.variables[name]; ok {
// 		e.variables[name] = val
// 		return nil
// 	}
// 	if e.enclosing != nil {
// 		return e.enclosing.Assign(name, val)
// 	}

// 	return fmt.Errorf("undefined variable '%s'", name)
// }

// func (e *Environment) Get(name string) (value.Value, error) {
// 	val, ok := e.variables[name]
// 	fmt.Printf("ENV LOOKING UP '%s' as type %v", name, val.Type)
// 	if ok {
// 		return val, nil
// 	}

// 	if e.enclosing != nil {
// 		return e.enclosing.Get(name)
// 	}

// 	return value.Null(), fmt.Errorf("undefined variable '%s'", name)
// }
