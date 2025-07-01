package symtable

import (
	"fmt"
	"strings"
	"sync"
)

type SymbolTable struct {
	Vars       map[string]Symbol
	Funcs      map[string]Symbol
	Types      map[string]Symbol
	TypeParams map[string]Symbol
	Parent     *SymbolTable
}

var genericTypeCache sync.Map

func NewSymbolTable(parent *SymbolTable) *SymbolTable {
	return &SymbolTable{
		Vars:       make(map[string]Symbol),
		Funcs:      make(map[string]Symbol),
		Types:      make(map[string]Symbol),
		TypeParams: make(map[string]Symbol),
		Parent:     parent,
	}
}

func InstantiationName(base string, args []*TypeSymbol) string {
	var argNames []string
	for _, t := range args {
		argNames = append(argNames, t.SymName)
	}
	return fmt.Sprintf("%s[%s]", base, strings.Join(argNames, ","))
}

func InstantiateGenericType(gen *TypeSymbol, typeArgs []*TypeSymbol) *TypeSymbol {
	if !gen.IsGeneric || len(gen.TypeParams) != len(typeArgs) {
		return gen
	}

	key := InstantiationName(gen.SymName, typeArgs)
	if val, ok := genericTypeCache.Load(key); ok {
		return val.(*TypeSymbol)
	}

	paramMap := map[string]*TypeSymbol{}
	for i, pname := range gen.TypeParams {
		paramMap[pname] = typeArgs[i]
	}

	newFields := map[string]*TypeSymbol{}
	for fname, ftype := range gen.Fields {
		newFields[fname] = substituteTypeParams(ftype, paramMap)
	}

	inst := &TypeSymbol{
		SymName:    key,
		SymKind:    gen.SymKind,
		TypeParams: nil,
		TypeArgs:   typeArgs,
		Fields:     newFields,
		Methods:    nil,
		IsGeneric:  false,
		Origin:     gen,
	}
	genericTypeCache.Store(key, inst)
	return inst
}

func substituteTypeParams(typ *TypeSymbol, paramMap map[string]*TypeSymbol) *TypeSymbol {
	if concrete, ok := paramMap[typ.SymName]; ok {
		return concrete
	}

	if typ.TypeArgs != nil && typ.Origin != nil {
		newArgs := make([]*TypeSymbol, len(typ.TypeArgs))
		for i, arg := range typ.TypeArgs {
			newArgs[i] = substituteTypeParams(arg, paramMap)
		}
		return InstantiateGenericType(typ.Origin, newArgs)
	}
	return typ
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
		// fmt.Printf("DEBUG: DefineValue Skipped duplicate: %q kind=%v\n", sym.Name(), sym.Kind())
		return fmt.Errorf("%s already defined in this scope", kindString(sym.Kind()))
	}
	ns[sym.Name()] = sym
	// fmt.Printf("DEBUG: DefineValue ADDED: %q kind=%v\n", sym.Name(), sym.Kind())
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
