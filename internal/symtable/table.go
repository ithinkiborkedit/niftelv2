package symtable

import "fmt"

type SymbolTable struct {
	values map[string]Symbol
	types  map[string]Symbol
	parent *SymbolTable
}

func NewSymbolTable(parent *SymbolTable) *SymbolTable {
	return &SymbolTable{
		values: make(map[string]Symbol),
		types:  make(map[string]Symbol),
		parent: parent,
	}
}

func (t *SymbolTable) DefineValue(sym Symbol) error {
	name := sym.Name()
	if _, exists := t.values[name]; exists {
		return fmt.Errorf("value '%s' already defined in this scope", name)
	}
	t.values[name] = sym
	return nil
}

func (t *SymbolTable) DefineType(sym Symbol) error {
	name := sym.Name()
	if _, exists := t.types[name]; exists {
		return fmt.Errorf("type '%s' already defined in this scope", name)
	}
	t.types[name] = sym
	return nil
}

func (t *SymbolTable) LookupValue(name string) (Symbol, bool) {
	sym, ok := t.values[name]
	if ok {
		return sym, true
	}
	if t.parent != nil {
		return t.parent.LookupValue(name)
	}
	return nil, false
}

func (t *SymbolTable) LookupType(name string) (Symbol, bool) {
	sym, ok := t.types[name]
	if ok {
		return sym, true
	}
	if t.parent != nil {
		return t.parent.LookupValue(name)
	}
	return nil, false
}

func (t *SymbolTable) HasLocalValue(name string) bool {
	_, ok := t.values[name]
	return ok
}

func (t *SymbolTable) HasLocalType(name string) bool {
	_, ok := t.types[name]
	return ok
}
