package symtable

import "fmt"

type SymbolTable struct {
	Vars       map[string]Symbol
	Funcs      map[string]Symbol
	Types      map[string]Symbol
	TypeParams map[string]Symbol
	Parent     *SymbolTable
}

func NewSymbolTable(parent *SymbolTable) *SymbolTable {
	return &SymbolTable{
		Vars:       make(map[string]Symbol),
		Funcs:      make(map[string]Symbol),
		Types:      make(map[string]Symbol),
		TypeParams: make(map[string]Symbol),
		Parent:     parent,
	}
}

func (s *SymbolTable) namespace(kind SymbolKind) map[string]Symbol {
	switch kind {
	case SymbolVar:
		return s.Vars
	case SymbolFuncs:
		return s.Funcs
	case SymbolTypes:
		return s.Types
	case SymbolTypeParams:
		return s.TypeParams
	default:
		panic("unknown symbol kind")
	}
}
func (s *SymbolTable) DefineValue(sym Symbol) error {
	ns := s.namespace(sym.Kind())
	if _, exists := ns[sym.Name()]; exists {
		return fmt.Errorf("%s already defined in this scope", kindString(sym.Kind()))
	}
	ns[sym.Name()] = sym
	return nil
}

func (s *SymbolTable) Lookup(kind SymbolKind, name string) (Symbol, bool) {
	ns := s.namespace(kind)
	if sym, ok := ns[name]; ok {
		return sym, true
	}
	if s.Parent != nil {
		return s.Parent.Lookup(kind, name)
	}
	return nil, false
}

func (s *SymbolTable) HasLocal(kind SymbolKind, name string) bool {
	ns := s.namespace(kind)
	_, ok := ns[name]
	return ok
}

func kindString(kind SymbolKind) string {
	switch kind {
	case SymbolVar:
		return "variable"
	case SymbolFuncs:
		return "function"
	case SymbolTypes:
		return "type"
	case SymbolTypeParams:
		return "typeparam"
	default:
		return "unknown"
	}
}
